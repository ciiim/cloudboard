package bptree

import "unsafe"

const (
	MaxRecordPerPage = (PageSize - unsafe.Sizeof(byte(0)) - unsafe.Sizeof(uint64(0)) - unsafe.Sizeof(uintptr(0))) / 16
)

// B+Tree 中的记录
type record struct {
	key   uint64
	value uint64
}

type nodeFlag byte

const (
	leafFlag nodeFlag = 1 << 0
	rootFlag nodeFlag = 1 << 1
)

func (f nodeFlag) isLeaf() bool {
	return f&leafFlag == 1
}

func (f nodeFlag) isRoot() bool {
	return f&rootFlag == 1
}

type data struct {
	gPID   globalPageID
	length int64 // 数据长度
}

type dataPage struct {
	flags nodeFlag
}

// B+Tree 中的 Node
// 4096 Bytes
type node struct {
	flags     nodeFlag
	recordNum uint64
	nextPage  globalPageID
	record    [MaxRecordPerPage]record // 最多存储 254 条记录
	_         uint64                   // 内存对齐
}

type BPTree struct {
	degree int
	root   *node
}
