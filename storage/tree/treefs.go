package tree

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ciiim/cloudborad/storage/types"
)

type DirEntry = fs.DirEntry

type TreeFileSystem struct {
	mu         sync.RWMutex
	usingSpace map[string]*Space

	config *Config
}

var _ ITreeFileSystem = (*TreeFileSystem)(nil)

const (
	DIR_PERFIX = "__DIR__"
	STAT_FILE  = "__STAT__"
	BASE_DIR   = "__BASE__"
)

func NewTreeFileSystem(config *Config) *TreeFileSystem {
	err := os.MkdirAll(config.RootPath, os.ModePerm)
	if err != nil {
		return nil
	}
	t := &TreeFileSystem{
		config: config,

		usingSpace: make(map[string]*Space),
	}
	return t
}

func (t *TreeFileSystem) NewLocalSpace(space string, cap types.Byte) error {

	if _, err := os.Stat(filepath.Join(t.config.RootPath, space, BASE_DIR)); err == nil {
		return ErrSpaceExist
	}

	err := os.Mkdir(filepath.Join(t.config.RootPath, space), 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Join(t.config.RootPath, space, BASE_DIR), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(t.config.RootPath, space, STAT_FILE))
	if err != nil {
		return err
	}
	defer file.Close()
	return err
}

func (t *TreeFileSystem) AllSpaces() []SpaceInfo {
	var spaces []SpaceInfo
	filepath.Walk(t.config.RootPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == BASE_DIR {
				space := filepath.Dir(path)
				space = filepath.Base(space)
				spaces = append(spaces, SpaceInfo{
					SpaceName: space,
				})
			}
		}
		return nil
	})
	return spaces

}

func (t *TreeFileSystem) GetLocalSpace(space string) *Space {
	s := func() *Space {
		t.mu.RLock()
		defer t.mu.RUnlock()
		if s, ok := t.usingSpace[space]; ok {
			return s
		}
		return nil
	}()
	if s != nil {
		return s
	}
	_, err := os.Stat(filepath.Join(t.config.RootPath, space, STAT_FILE))
	if err != nil {
		log.Println("[Space] Missing stat file: ", space)
		return nil
	}
	s = &Space{
		root:     t.config.RootPath,
		spaceKey: space,
		base:     BASE_DIR,
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	t.usingSpace[space] = s

	return s
}

func (t *TreeFileSystem) DeleteLocalSpace(space string) error {
	if _, err := os.Stat(filepath.Join(t.config.RootPath, space, BASE_DIR)); err != nil {
		return ErrSpaceNotFound
	}
	return os.RemoveAll(filepath.Join(t.config.RootPath, space))
}

func (t *TreeFileSystem) MakeDir(space, base, name string) error {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return ErrSpaceNotFound
	}
	return sp.makeDir(base, name)

}

func (t *TreeFileSystem) RenameDir(space, base, name, newName string) error {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return ErrSpaceNotFound
	}
	return sp.renameDir(base, name, newName)

}

func (t *TreeFileSystem) DeleteDir(space, base, name string) error {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return ErrSpaceNotFound
	}
	return sp.deleteDir(base, name)

}

func (t *TreeFileSystem) GetDirSub(space, base, name string) ([]*SubInfo, error) {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return nil, ErrSpaceNotFound
	}
	dirPath, des, err := sp.getDir(base, name)
	if err != nil {
		return nil, err
	}
	return DirEntryToSubInfo(dirPath, des), nil

}

func (t *TreeFileSystem) GetMetadata(space, base, name string) ([]byte, error) {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return nil, ErrSpaceNotFound
	}
	return sp.getMetadata(base, name)

}

func (t *TreeFileSystem) PutMetadata(space, base, name string, fileHash []byte, metadata []byte) error {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return ErrSpaceNotFound
	}

	return sp.storeMetaData(base, name, metadata)

}
func (t *TreeFileSystem) DeleteMetadata(space, base, name string) error {
	sp := t.GetLocalSpace(space)
	if sp == nil {
		return ErrSpaceNotFound
	}
	return sp.deleteMetaData(base, name)
}

func (t *TreeFileSystem) Close() error {
	return nil
}
