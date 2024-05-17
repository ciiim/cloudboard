package model

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"time"

	"github.com/ciiim/cloudborad/storage/tree"
	"github.com/ciiim/cloudborad/storage/types"
)

const (
	STAT_PASSWORD = "password"
)

var reservedStat = []string{STAT_PASSWORD}

// 用户目录
type UserDir struct {
	FileNums int         `json:"file_nums"`
	Files    []*UserFile `json:"files"`
}

type UserFile struct {
	IsDir      bool   `json:"is_dir"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	ModTime    int64  `json:"mod_time"`
	CreateTime int64  `json:"create_time"`
}

func handleBaseDir(base string) string {
	if base == "" || base == "/" {
		base = tree.BASE_DIR
	}
	return base
}

func (r *RingModel) NewSpace(space string, cap types.Byte) error {
	if cap < 0 {
		// FIXME: return error
		return nil
	}
	return r.FrontSystem.NewSpace(space, cap)
}

func (r *RingModel) GetAllSpaces() []tree.SpaceInfo {
	return r.FrontSystem.AllSpaces()
}

func (r *RingModel) IsReservedStat(space, name string) bool {
	for _, v := range reservedStat {
		if v == name {
			return true
		}
	}
	return false
}

func (r *RingModel) GetSpaceStat(space, key string) (*tree.SpaceStatElement, error) {
	if r.IsReservedStat(space, key) {
		return nil, fmt.Errorf("reserved stat")
	}
	return r.FrontSystem.GetSpaceStat(space, key)
}

func (r *RingModel) SetSpaceStat(space, key, value string) error {
	if r.IsReservedStat(space, key) {
		return fmt.Errorf("reserved stat")
	}
	return r.FrontSystem.SetSpaceStat(space, tree.NewSpaceStatElement(key, value))
}

func (r *RingModel) SetSpacePassword(space, password string) error {
	if err := r.FrontSystem.SetSpaceStat(space, tree.NewSpaceStatElement(STAT_PASSWORD, password)); err != nil {
		return err
	}
	return nil
}

func (r *RingModel) CheckSpacePassword(space, password string) bool {
	stat, err := r.FrontSystem.GetSpaceStat(space, STAT_PASSWORD)
	if err != nil {
		log.Println(err)
		return false
	}
	return stat.Value() == password
}

func (r *RingModel) CheckSpaceToken(space, token string) bool {
	pwd, err := r.FrontSystem.GetSpaceStat(space, "password")
	if err != nil {
		return false
	}
	return DecryptToken(token) == pwd.Value()
}

func (r *RingModel) MakeDir(space, base, name string) error {
	return r.FrontSystem.MakeDir(space, base, name)
}

func (r *RingModel) RenameDir(space, base, name, newName string) error {
	return r.FrontSystem.RenameDir(space, base, name, newName)
}

func (r *RingModel) DeleteDir(space, base, name string) error {
	return r.FrontSystem.DeleteDir(space, base, name)
}

func (r *RingModel) DirInSpace(space string, base, dir string) (*UserDir, error) {
	base = handleBaseDir(base)
	subInfo, err := r.FrontSystem.GetDirSub(space, base, dir)
	if err != nil {
		return nil, err
	}
	userDir := &UserDir{}
	userDir.FileNums = len(subInfo)
	userDir.Files = make([]*UserFile, userDir.FileNums)
	for i, sub := range subInfo {
		userDir.Files[i] = &UserFile{
			IsDir:      sub.IsDir,
			FileName:   sub.Name,
			FileSize:   sub.Size,
			ModTime:    sub.ModTime.Unix(),
			CreateTime: sub.CreateTime.Unix(),
		}
	}
	return userDir, nil
}

func (r *RingModel) PutFile(space, base, name string, fileHash []byte, fileSize int64, file multipart.File) error {
	base = handleBaseDir(base)
	chunkMaxSize := r.StorageSystem.Config().HCSConfig.ChunkMaxSize
	if fileSize > chunkMaxSize {
		return r.putFileSplit(space, base, name, chunkMaxSize, fileHash, fileSize, file)
	}
	return r.putFile(space, base, name, fileHash, fileSize, file)
}

// 直接存储文件
func (r *RingModel) putFile(space, base, name string, fileHash []byte, fileSize int64, file multipart.File) error {

	chunks := make([]*tree.FileChunk, 1)

	//创建FileChunk
	chunks[0] = tree.NewFileChunk(fileSize, fileHash)

	//存储元数据
	metadata := tree.NewMetaData(name, fileHash, time.Now(), chunks)
	metadataBytes, err := tree.MarshalMetaData(metadata)
	if err != nil {
		return err
	}
	if err = r.FrontSystem.PutMetadata(space, base, name, fileHash, metadataBytes); err != nil {
		return err
	}

	//存储文件
	return r.StorageSystem.StoreReader(fileHash, hex.EncodeToString(fileHash), fileSize, file, nil)
}

// 分片存储文件
func (r *RingModel) putFileSplit(space, base, name string, chunkSize int64, fileHash []byte, fileSize int64, file multipart.File) error {
	// 计算分片数量
	chunkNum := int(math.Ceil(float64(fileSize) / float64(chunkSize)))
	//创建分片Reader
	chunkReaders := make([]*io.SectionReader, chunkNum)
	chunks := make([]*tree.FileChunk, chunkNum)
	for i := 0; i < chunkNum-1; i++ {
		chunkReaders[i] = io.NewSectionReader(file, int64(i)*chunkSize, chunkSize)
	}

	//最后一个分片Reader
	chunkReaders[chunkNum-1] = io.NewSectionReader(file, int64(chunkNum-1)*chunkSize, fileSize-int64(chunkNum-1)*chunkSize)

	//计算每个分片的hash
	if err := func() error {
		chunkBuffer := r.chunkPool.Get()
		defer r.chunkPool.Put(chunkBuffer)
		for i, chunkReader := range chunkReaders {
			chunkBuffer.Reset()
			_, err := io.Copy(chunkBuffer, chunkReader)
			if err != nil {
				return err
			}
			chunks[i] = tree.NewFileChunk(
				chunkReader.Size(),
				chunkBuffer.Hash(r.StorageSystem.Config().HCSConfig.HashFn),
			)
			if _, err := chunkReaders[i].Seek(0, io.SeekStart); err != nil {
				return err
			}
		}
		return nil
	}(); err != nil {
		return err
	}

	//存储元数据
	metadata := tree.NewMetaData(name, fileHash, time.Now(), chunks)
	metadataBytes, err := tree.MarshalMetaData(metadata)
	if err != nil {
		return err
	}

	if err = r.FrontSystem.PutMetadata(space, base, name, fileHash, metadataBytes); err != nil {
		return err
	}

	//存储分片
	for i, chunkReader := range chunkReaders {
		if err := r.StorageSystem.StoreReader(
			chunks[i].Hash,
			hex.EncodeToString(chunks[i].Hash),
			chunks[i].Size,
			chunkReader,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}

type fileReader struct {
	io.ReadCloser
	FileSize int64
	FileName string
}

// 下载文件
func (r *RingModel) GetFile(space, base, name string) (*fileReader, error) {
	base = handleBaseDir(base)
	metadataBytes, err := r.FrontSystem.GetMetadata(space, base, name)
	if err != nil {
		fmt.Printf("GetMetadata: %v\n", err)
		return nil, err
	}
	var metadata tree.Metadata
	if err = tree.UnmarshalMetaData(metadataBytes, &metadata); err != nil {
		fmt.Printf("UnmarshalMetaData: %v\n", err)
		return nil, err
	}

	var (
		chunkClosers []io.Closer
		chunkReaders []io.Reader
		multiReader  io.Reader
	)

	for _, v := range metadata.Chunks {
		chunk, err := r.StorageSystem.Get(v.Hash)
		if err != nil {
			fmt.Printf("GetFile: %v\n", err)
			return nil, err
		}
		chunkReaders = append(chunkReaders, chunk)
		chunkClosers = append(chunkClosers, chunk)
	}

	multiReader = io.MultiReader(chunkReaders...)

	return &fileReader{
		ReadCloser: multiReadCloser(multiReader, func() error {
			for _, v := range chunkClosers {
				if err := v.Close(); err != nil {
					return err
				}
			}
			return nil
		}),
		FileSize: metadata.Size,
		FileName: metadata.Filename,
	}, nil
}

func (r *RingModel) DeleteFile(space, base, name string) error {
	base = handleBaseDir(base)
	metadataBytes, err := r.FrontSystem.GetMetadata(space, base, name)
	if err != nil {
		return err
	}

	var metadata tree.Metadata

	if err = tree.UnmarshalMetaData(metadataBytes, &metadata); err != nil {
		return err
	}

	for _, v := range metadata.Chunks {
		if err := r.StorageSystem.Delete(v.Hash); err != nil {
			return err
		}
		fmt.Printf("DeleteChunk: %v\n", hex.EncodeToString(v.Hash))
	}

	return r.FrontSystem.DeleteMetadata(space, base, name)
}
