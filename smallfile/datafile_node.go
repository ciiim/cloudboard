package smallfile

import (
	"encoding/binary"
	"io"
	"unsafe"
)

const (
	dataBodyLength = unsafe.Sizeof(GlobalPageID(0)) + unsafe.Sizeof(dataKey64Bytes{}) + unsafe.Sizeof(uint64(0)) + unsafe.Sizeof(GlobalPageID(0))
)

type dataKey64Bytes [64]byte

type dataHead struct {
	gPID           GlobalPageID    // offset:0
	key            *dataKey64Bytes // offset:8
	length         uint64          // 数据长度 offset:72
	nextDataPageID GlobalPageID    // 下一个数据页的ID offset:80 - 88
	// data   []byte          // 数据
}

const (
	offset1 = 0
	offset2 = 8
	offset3 = 72
	offset4 = 80
	offset5 = 88
)

type dataSpan struct {
	key         *dataKey64Bytes
	limitReader *io.LimitedReader
}

func parseDataHead(dataSpan *SpanInCache) *dataHead {
	d := &dataHead{}

	d.gPID = dataSpan.globalID
	d.key = (*dataKey64Bytes)(dataSpan.space.buf[offset1:offset2])
	d.length = binary.BigEndian.Uint64(dataSpan.space.buf[offset3:offset4])
	d.nextDataPageID = GlobalPageID(binary.BigEndian.Uint64(dataSpan.space.buf[offset4:offset5]))
	return d
}

func makeDataHead(pageID GlobalPageID, dataLen uint64, key *dataKey64Bytes) *dataHead {
	d := &dataHead{}
	d.gPID = pageID
	d.key = key
	d.length = dataLen
	d.nextDataPageID = 0 //FIXME: 以后实现
	return d
}

func dataHeadToBytes(d *dataHead) []byte {
	buf := make([]byte, dataBodyLength)
	binary.BigEndian.PutUint64(buf[offset1:offset2], uint64(d.gPID))
	copy(buf[offset2:offset3], d.key[:])
	binary.BigEndian.PutUint64(buf[offset3:offset4], d.length)
	binary.BigEndian.PutUint64(buf[offset4:offset5], uint64(d.nextDataPageID))
	return buf
}
