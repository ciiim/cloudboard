package smallfile

import (
	"bytes"
	"io"
)

const (
	dataFileName = "data.db"
)

type DataFile struct {
	info      GlobalPageID
	allocator *Allocator
}

func NewDataFile() *DataFile {
	return &DataFile{
		allocator: NewAllocator(dataFileName),
	}
}

func (df *DataFile) Put(key *dataKey64Bytes, data []byte) (GlobalPageID, error) {
	length := len(data) + int(dataBodyLength)
	dataPageID, err := df.allocator.Alloc(length)
	if err != nil {
		return 0, err
	}

	return dataPageID, nil
}

func (df *DataFile) Get(id GlobalPageID) (*dataSpan, error) {
	data, err := df.allocator.Get(id)
	if err != nil {
		return nil, err
	}
	dataHead := parseDataHead(data)
	dataBodyReader := bytes.NewReader(data.space.buf[offset5:])
	return &dataSpan{
		key: dataHead.key,
		limitReader: &io.LimitedReader{
			R: dataBodyReader,
			N: int64(dataHead.length),
		},
	}, nil
}

func (df *DataFile) Delete(id GlobalPageID) {
	df.allocator.Free(id)
}
