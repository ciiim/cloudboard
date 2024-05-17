// implement hash chunk system
package hashchunk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ciiim/cloudborad/database"
	"github.com/ciiim/cloudborad/storage/types"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrEmptyKey = errors.New("key is empty")
)

type HashChunkSystem struct {
	config *Config

	capacity *types.SafeInt64
	occupied *types.SafeInt64

	chunkStatDBName string

	levelDB *leveldb.DB

	rwMutex sync.RWMutex
}

type CalcChunkStoragePathFn = func(chunkStat *HashChunkInfo) string

type Hash func([]byte) []byte

// default calculate store path function
// format: chunkhash[0:3]/chunkhash[3:6]/chunkhash[6:9]
var DefaultCalcStorePathFn = func(hci *HashChunkInfo) string {
	path := fmt.Sprintf("%x/%x/%x", hci.ChunkHash[0:3], hci.ChunkHash[3:6], hci.ChunkHash[3:9])
	return path
}

var _ IHashChunkSystem = (*HashChunkSystem)(nil)
var _ IHashChunkInfo = (*HashChunkInfo)(nil)

func NewHashChunkSystem(config *Config) *HashChunkSystem {
	if err := os.MkdirAll(config.RootPath, os.ModePerm); err != nil {
		panic("mkdir error:" + err.Error())
	}
	hashDBName := "chunkinfo"
	db, err := database.NewLevelDB(filepath.Join(config.RootPath + "/" + hashDBName))
	if err != nil {
		panic("leveldb init error:" + err.Error())
	}

	hcs := &HashChunkSystem{
		capacity:        types.NewSafeInt64(),
		occupied:        types.NewSafeInt64(),
		chunkStatDBName: hashDBName,
		levelDB:         db,
		config:          config,
	}
	if config.CalcStoragePathFn == nil {
		hcs.config.CalcStoragePathFn = DefaultCalcStorePathFn
	}

	cap, ouppy, err := hcs.getCapAndOccupied()

	if err != nil {
		log.Println("New HCS at", config.RootPath)
		_ = hcs.storeCapAndOccupied(config.Capacity, 0)
		return hcs
	}
	log.Printf("Detect exist HCS at %s\n", config.RootPath)

	hcs.capacity.Store(cap)
	hcs.updateOccupied(ouppy)

	if config.Capacity < cap {
		log.Println("[BFS] capacity is less than exist HCS, use exist HCS's capacity.")
	}
	if config.Capacity > cap {
		log.Println("[BFS] capacity is more than exist HCS, use new capacity.")
		hcs.capacity.Store(config.Capacity)
	}
	return hcs
}

func (hcs *HashChunkSystem) CreateChunk(key []byte, chunkName string, size int64, extra *ExtraInfo) (io.WriteCloser, error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	}

	hcs.rwMutex.Lock()
	defer hcs.rwMutex.Unlock()

	// increase chunk counter
	// if chunk is exist, just increase counter
	// if chunk is not exist, create chunk info and store it
	_, err := hcs.increaseChunkCounter(key)
	if err == nil {
		return nil, nil
	}
	if err != nil && err != ErrChunkInfoNotFound && err != ErrChunkNotFound {
		return nil, err
	}

	fmt.Printf("CreateChunk: key: %x, chunkName: %s, size: %d\n", key, chunkName, size)

	hci := NewChunkInfo(chunkName, key, size)
	hci.SetPath(filepath.Join(hcs.config.RootPath, hcs.config.CalcStoragePathFn(hci)))
	info := NewInfo(hci, extra)
	if err := os.MkdirAll(hci.ChunkPath, os.ModePerm); err != nil {
		return nil, err
	}
	if err := hcs.storeInfo(key, info); err != nil {
		return nil, err
	}
	return hcs.createChunkWriter(hci)
}

func (hcs *HashChunkSystem) StoreBytes(key []byte, chunkName string, size int64, value []byte, extra *ExtraInfo) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if value == nil {
		return fmt.Errorf("value is nil")
	}

	hcs.rwMutex.Lock()
	defer hcs.rwMutex.Unlock()

	// increase chunk counter
	// if chunk is exist, just increase counter
	// if chunk is not exist, create chunk info and store it
	_, err := hcs.increaseChunkCounter(key)
	if err == nil {
		return nil
	}
	if err != nil && err != ErrChunkInfoNotFound {
		return err
	}

	valueLength := int64(len(value))

	if valueLength != size {
		return fmt.Errorf("value length is not equal to size")
	}

	//check capacity
	if err := hcs.CheckCapacity(valueLength); err != nil {
		return err
	}

	hci := NewChunkInfo(chunkName, key, size)
	hci.SetPath(filepath.Join(hcs.config.RootPath, hcs.config.CalcStoragePathFn(hci)))

	info := NewInfo(hci, extra)
	//make dir
	if err := os.MkdirAll(hci.ChunkPath, os.ModePerm); err != nil {
		return err
	}
	if err := hcs.storeInfo(key, info); err != nil {
		return err
	}
	if err := hcs.storeChunkBytes(hci, value); err != nil {
		return err
	}

	//update Occupied
	hcs.updateOccupied(hcs.occupied.Load() + hci.ChunkSize)

	return nil
}

func (hcs *HashChunkSystem) StoreReader(key []byte, chunkName string, size int64, v io.Reader, extra *ExtraInfo) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}
	if v == nil {
		return fmt.Errorf("value is nil")
	}

	hcs.rwMutex.Lock()
	defer hcs.rwMutex.Unlock()

	// increase chunk counter
	// if chunk is exist, just increase counter
	// if chunk is not exist, create chunk info and store it
	_, err := hcs.increaseChunkCounter(key)
	if err == nil {
		return nil
	}

	if err != nil && err != ErrChunkInfoNotFound {
		return err
	}

	//check capacity
	hci := NewChunkInfo(chunkName, key, size)
	hci.SetPath(filepath.Join(hcs.config.RootPath, hcs.config.CalcStoragePathFn(hci)))
	info := NewInfo(hci, extra)

	if err := os.MkdirAll(hci.ChunkPath, os.ModePerm); err != nil {
		return err
	}
	if err := hcs.storeInfo(key, info); err != nil {
		return err
	}
	if err := hcs.storeChunkReader(hci, v); err != nil {
		return err
	}

	//update Occupied
	hcs.updateOccupied(hcs.occupied.Load() + hci.ChunkSize)

	return nil
}

func (hcs *HashChunkSystem) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, ErrEmptyKey
	}

	hcs.rwMutex.RLock()
	defer hcs.rwMutex.RUnlock()

	return hcs.isExist(key), nil

}

func (hcs *HashChunkSystem) Get(key []byte) (*HashChunk, error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	}

	hcs.rwMutex.RLock()
	defer hcs.rwMutex.RUnlock()

	info, err := hcs.getInfo(key)
	if err != nil {
		return nil, err
	}
	file, err := hcs.getChunk(info)
	return &HashChunk{
		ReadSeekCloser: file,
		info:           info,
	}, err
}

func (hcs *HashChunkSystem) Delete(key []byte) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}

	hcs.rwMutex.Lock()
	defer hcs.rwMutex.Unlock()

	// decrease chunk counter
	nowCounter, err := hcs.decreaseChunkCounter(key)

	// still have reference
	if err == nil && nowCounter != 0 {
		return nil
	}

	if err != nil {
		return err
	}

	info, err := hcs.getInfo(key)
	if err != nil {
		return err
	}
	// TODO: check occupied
	// if hcs.occupied.Load()-info.ChunkInfo.ChunkSize < 0 {
	// 	return fmt.Errorf("[Delete Chunk Error] Occupied is 0")
	// }
	if err := hcs.DeleteInfo(key); err != nil {
		return err
	}
	if err := hcs.deleteChunk(info); err != nil {
		return err
	}
	//update Occupied
	hcs.updateOccupied(hcs.occupied.Load() - info.ChunkInfo.ChunkSize)
	return nil
}

func (hcs *HashChunkSystem) GetInfo(key []byte) (*Info, error) {
	if len(key) == 0 {
		return nil, ErrEmptyKey
	}

	hcs.rwMutex.RLock()
	defer hcs.rwMutex.RUnlock()

	return hcs.getInfo(key)
}

func (hcs *HashChunkSystem) UpdateInfo(key []byte, updateInfoFn func(info *Info)) error {
	if len(key) == 0 {
		return ErrEmptyKey
	}

	hcs.rwMutex.Lock()
	defer hcs.rwMutex.Unlock()

	info, err := hcs.getInfo(key)
	if err != nil {
		return err
	}

	updateInfoFn(info)

	return hcs.storeInfo(key, info)
}

func (hcs *HashChunkSystem) Opt(opt any) any {
	return nil
}

func (hcs *HashChunkSystem) isExist(key []byte) bool {
	_, err := hcs.getInfo(key)
	return err == nil
}

func (hcs *HashChunkSystem) Cap() int64 {
	return hcs.capacity.Load()
}

// unit can be "B", "KB", "MB", "GB" or just leave it blank
func (hcs *HashChunkSystem) Occupied(unit ...string) float64 {
	if len(unit) == 0 {
		return float64(hcs.occupied.Load())
	}
	switch unit[0] {
	case "B":
		return float64(hcs.occupied.Load())
	case "KB":
		return float64(hcs.occupied.Load()) / 1024
	case "MB":
		return float64(hcs.occupied.Load()) / 1024 / 1024
	case "GB":
		return float64(hcs.occupied.Load()) / 1024 / 1024 / 1024
	default:
		return float64(hcs.occupied.Load())
	}
}

func (hcs *HashChunkSystem) createChunkWriter(chunkInfo *HashChunkInfo) (io.WriteCloser, error) {
	if hcs.config.CalcStoragePathFn == nil {
		return nil, fmt.Errorf("CalcChunkStoragePathFn is nil")
	}
	chunkFile, err := os.Create(filepath.Join(chunkInfo.ChunkPath, chunkInfo.ChunkName))
	if err != nil {
		return nil, fmt.Errorf("open file %s error: %w", chunkInfo.ChunkPath+"/"+chunkInfo.ChunkName, err)
	}
	chunkwc := warpHashChunkWriteCloser(chunkFile)
	return chunkwc, nil
}

func (hcs *HashChunkSystem) storeChunkBytes(hcStat *HashChunkInfo, value []byte) error {
	file, err := os.Create(filepath.Join(hcStat.ChunkPath, hcStat.ChunkName))
	if err != nil {
		return fmt.Errorf("open file %s error: %s", hcStat.ChunkPath+"/"+hcStat.ChunkName, err)
	}
	defer file.Close()

	_, err = file.Write(value)
	return err
}

func (hcs *HashChunkSystem) storeChunkReader(key *HashChunkInfo, reader io.Reader) error {
	if hcs.config.CalcStoragePathFn == nil {
		return fmt.Errorf("CalcChunkStoragePathFn is nil")
	}
	file, err := os.Create(key.ChunkPath + "/" + key.ChunkName)
	if err != nil {
		return fmt.Errorf("open file %s error: %s", key.ChunkPath+"/"+key.ChunkName, err)
	}
	_, err = file.ReadFrom(reader)
	return err
}

func (hcs *HashChunkSystem) getChunk(info *Info) (io.ReadSeekCloser, error) {
	file, err := os.Open(filepath.Join(info.ChunkInfo.ChunkPath, info.ChunkInfo.ChunkName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrChunkNotFound
		}
		return nil, err
	}
	return file, nil
}

func (hcs *HashChunkSystem) deleteChunk(info *Info) error {
	fullPath := filepath.Join(info.ChunkInfo.ChunkPath, info.ChunkInfo.ChunkName)
	err := os.Remove(fullPath)
	return err
}

func (hcs *HashChunkSystem) getInfo(hashSum []byte) (*Info, error) {
	infoBytes, err := hcs.levelDB.Get(hashSum, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrChunkInfoNotFound
		}
		return nil, err
	}
	var info Info
	err = json.Unmarshal(infoBytes, &info)
	return &info, err
}

func (hcs *HashChunkSystem) storeInfo(hashSum []byte, info *Info) error {
	res, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = hcs.levelDB.Put(hashSum, res, nil)
	return err
}

func (hcs *HashChunkSystem) DeleteInfo(hashSum []byte) error {
	return hcs.levelDB.Delete(hashSum, nil)
}

func (hcs *HashChunkSystem) increaseChunkCounter(key []byte) (nowCounter int64, err error) {
	// get chunk info
	info, err := hcs.getInfo(key)
	if err != nil {
		return 0, err
	}

	//检查chunk是否存在
	if _, err = os.Stat(filepath.Join(info.ChunkInfo.ChunkPath, info.ChunkInfo.ChunkName)); err != nil {
		if os.IsNotExist(err) {
			return 0, ErrChunkNotFound
		}
	}

	info.ChunkInfo.ChunkCount++

	// store chunk info
	return info.ChunkInfo.ChunkCount, hcs.storeInfo(key, info)
}

func (hcs *HashChunkSystem) decreaseChunkCounter(key []byte) (nowCounter int64, err error) {
	// get chunk info
	info, err := hcs.getInfo(key)
	if err != nil {
		return 0, err
	}
	info.ChunkInfo.ChunkCount--
	if info.ChunkInfo.ChunkCount <= 0 {
		return 0, nil
	}

	// store chunk info
	return info.ChunkInfo.ChunkCount, hcs.storeInfo(key, info)

}

func (hcs *HashChunkSystem) storeCapAndOccupied(capacity, occupied int64) error {
	var capAndOccupied struct {
		Capacity int64 `json:"capacity"`
		Occupied int64 `json:"occupied"`
	}
	capAndOccupied.Capacity = capacity
	capAndOccupied.Occupied = occupied
	res, err := json.Marshal(capAndOccupied)
	if err != nil {
		return err
	}
	err = hcs.levelDB.Put([]byte("cap_and_occupied"), res, nil)
	return err
}

func (hcs *HashChunkSystem) getCapAndOccupied() (int64, int64, error) {

	res, err := hcs.levelDB.Get([]byte("cap_and_occupied"), nil)
	if err != nil {
		return 0, 0, err
	}
	var capAndOccupied struct {
		Capacity int64 `json:"capacity"`
		Occupied int64 `json:"occupied"`
	}
	err = json.Unmarshal(res, &capAndOccupied)
	if err != nil {
		return 0, 0, err
	}
	return capAndOccupied.Capacity, capAndOccupied.Occupied, nil
}

func (hcs *HashChunkSystem) Close() error {
	if hcs == nil {
		panic("HashFileSystem is nil")
	}
	log.Println("HashFileSystem Closing.")

	//save cap and ouppy
	if err := hcs.storeCapAndOccupied(hcs.capacity.Load(), hcs.occupied.Load()); err != nil {
		log.Println("Save filesystem error:", err)
	}
	return hcs.levelDB.Close()
}

func (hcs *HashChunkSystem) updateOccupied(occupied int64) {
	hcs.occupied.Store(occupied)
}

func (hcs *HashChunkSystem) CheckCapacity(delta int64) error {
	if hcs.occupied.Load()+delta > hcs.capacity.Load() {
		return ErrFull
	}
	return nil
}

func (hcs *HashChunkSystem) Config() *Config {
	return hcs.config
}