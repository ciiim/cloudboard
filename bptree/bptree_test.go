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

func TestSize(t *testing.T) {

	t.Log(unsafe.Sizeof(testSpan{}))
	t.Log(unsafe.Alignof(testSpan{}))
}
