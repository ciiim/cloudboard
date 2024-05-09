package bptree

import (
	"testing"
	"unsafe"
)

type testSpan struct {

	/*
		0bit: 1bit flag

		1-12bit: 12bits span length

		13-31bit: 19bits next free page index
	*/
	span uint32
}

type testNode struct {
	flag      byte
	recordNum uint8
	next      uint64
	records   [55]testRecord2
}

type testRecord1 struct {
	value  uint64
	keyLen uint8
	key    [64]byte
}

type testRecord2 struct {
	keyLen uint8
	key    [64]byte
	value  uint64
}

func TestSize(t *testing.T) {

	t.Log(unsafe.Sizeof(testRecord1{}))
	t.Log(unsafe.Alignof(testRecord1{}))
	t.Log(unsafe.Sizeof(testRecord2{}))
	t.Log(unsafe.Alignof(testRecord2{}))
	t.Log(unsafe.Sizeof(testNode{}))
	t.Log(unsafe.Alignof(testNode{}))
}
