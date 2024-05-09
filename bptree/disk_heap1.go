package bptree

import (
	"encoding/binary"
	"errors"
	"sync"
	"syscall"
)

const (
	PageSize = 1 << 12 //4KB

	HeaderMetadataByteNum = PageSize * 256 // 1MB

	BusyPageOffset = PageSize * 256 // 1MB

	BusyPageByteNum = PageSize * 256 // 1MB

	ChunkSpaceOffset = PageSize * 512 // 2MB

	MaxChunkNum = 1024 * 1024 / 4 // 最多chunk数 1MB / 4Byte

	MaxChunkSize = 1 << 31 // 2GB

	MaxPagePerChunk = MaxChunkSize / PageSize

	// 一个Chunk中系统使用的Page数，包含PageBitmap和SpanMap占用的page
	SysUsedPagePerChunk = SpanMapLengthPerChunk / PageSize // 1024 * 4KB

	UserPageNo = SysUsedPagePerChunk

	// 一个Chunk中用户可用的Page数
	MaxUserPagePerChunk = MaxPagePerChunk - SysUsedPagePerChunk
)

const (
	SpanInfoSize = 1 << 3 // 8Byte

	SpanMapLengthPerChunk = MaxPagePerChunk * SpanInfoSize

	MaxSpanPageNum = MaxUserPagePerChunk
)

var (
	ErrInvaildParam = errors.New("invaild param")

	ErrMaxChunkNum = errors.New("max chunk num")

	ErrMaxPagePerSpan = errors.New("max page per span")

	ErrMaxPageNum = errors.New("max page num")
)

// 全局pageID
type GlobalPageID uint64

func (g GlobalPageID) chunkID() chunkID {
	return chunkID(g / MaxPagePerChunk)
}

func (g GlobalPageID) toLocal() localPageID {
	return localPageID(g % MaxPagePerChunk)
}

// 本地pageID
type localPageID uint32

func (l localPageID) toGlobal(chunkID chunkID) GlobalPageID {
	return GlobalPageID(uint64(chunkID)*MaxPagePerChunk + uint64(l))
}

type chunkID = uint64

type headerMetadata struct {
	space space
}

func (h *headerMetadata) GetNextChunk() chunkID {
	return chunkID(binary.LittleEndian.Uint64(h.space.buf[0:8]))
}

func (h *headerMetadata) SetNextChunk(chunkID chunkID) {
	binary.LittleEndian.PutUint64(h.space.buf[0:8], uint64(chunkID))
}

type headerBusyPage struct {
	// chunks [MaxChunkNum]int32

	lock sync.RWMutex

	// mmap后得到该header的buf
	buf space
}

// 获取该chunk的busy page数
func (h *headerBusyPage) getChunk(index chunkID) int32 {

	if (int(index)+1)*int(h.buf.step) > len(h.buf.buf) {
		return -1
	}
	h.lock.RLock()
	defer h.lock.RUnlock()

	return int32(binary.LittleEndian.Uint32(h.buf.buf[int(index)*int(h.buf.step) : (int(index)+1)*int(h.buf.step)]))
}

func (h *headerBusyPage) addChunkBusyPages(index chunkID, delta int32) {

	if (int(index)+1)*int(h.buf.step) > len(h.buf.buf) {
		println("update chunk busy page failed1")
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()

	old := int32(binary.LittleEndian.Uint32(h.buf.buf[int(index)*int(h.buf.step) : (int(index)+1)*int(h.buf.step)]))

	if old+delta < SysUsedPagePerChunk || old+delta > MaxPagePerChunk {
		println("update chunk busy page failed2")
		return
	}

	binary.LittleEndian.PutUint32(h.buf.buf[index*uint64(h.buf.step):(index+1)*uint64(h.buf.step)], uint32(old+delta))

}

func (h *headerBusyPage) setChunkBusyPages(index chunkID, busyPage uint32) {

	h.lock.Lock()
	defer h.lock.Unlock()

	binary.LittleEndian.PutUint32(h.buf.buf[index*uint64(h.buf.step):(index+1)*uint64(h.buf.step)], busyPage)
}

type space struct {
	buf []byte

	// 每个数据的大小 unit:Byte
	step uint
}

func (s *space) Release() error {
	if s.buf == nil {
		return nil
	}
	return syscall.Munmap(s.buf)
}
