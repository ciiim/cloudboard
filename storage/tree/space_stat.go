package tree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

/*
SpaceStat 文件格式
一个元素一行，元素格式为 key:value
*/

type SpaceStatElement struct {
	key   string
	value string
}

func NewSpaceStatElement(key string, value string) *SpaceStatElement {
	return &SpaceStatElement{
		key:   key,
		value: value,
	}
}

func (e *SpaceStatElement) Key() string {
	return e.key
}

func (e *SpaceStatElement) Value() string {
	return e.value
}

func (s *Space) getStatPath() string {
	return filepath.Join(s.root, s.spaceKey, STAT_FILE)
}

func (s *Space) SetStatElement(e *SpaceStatElement) error {
	path := s.getStatPath()

	spaceStat, err := os.Open(path)
	if err != nil {
		return err
	}
	defer spaceStat.Close()

	tmpPath := path + ".tmp"
	tmpFile, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	var covered bool

	// 逐行读取文件，找到对应的key，修改value后写入临时文件
	var k, v string
	for {
		_, err := fmt.Fscanf(spaceStat, "%s:%s\n", &k, &v)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if e.key == k {
			covered = true
			// 如果value为空，则不写入临时文件
			if e.value == "" {
				continue
			}
			v = e.value
		}
		_, err = fmt.Fprintf(tmpFile, "%s %s\n", k, v)
		if err != nil {
			return err
		}
	}

	// 未被覆盖的值
	if !covered {
		_, err = fmt.Fprintf(tmpFile, "%s %s\n", e.key, e.value)
		if err != nil {
			return err
		}
	}

	// 删除原文件，重命名临时文件
	if err = os.Remove(path); err != nil {
		return err
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return err
	}
	return nil
}

func (s *Space) GetStatElement(key string) (*SpaceStatElement, error) {
	path := s.getStatPath()

	spaceStat, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer spaceStat.Close()
	// 逐行读取文件，找到对应的key
	var k, v string
	for {
		_, err := fmt.Fscanf(spaceStat, "%s %s\n", &k, &v)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if key == k {
			return NewSpaceStatElement(key, v), nil
		}
	}
	return nil, nil
}
