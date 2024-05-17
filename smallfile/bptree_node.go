package smallfile

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

const (
	RecordSize       = 73 // 1 keyLen + 64 key + 8 value
	MaxRecordPerPage = (PageSize - unsafe.Sizeof(byte(0)) - unsafe.Sizeof(uint8(0)) - unsafe.Sizeof(uint64(0))) / uintptr(RecordSize)
)

// uint64
type recordUint64 [8]byte

func (ru *recordUint64) toUint64() uint64 {
	return binary.LittleEndian.Uint64(ru[:])
}

type record64Bytes [64]byte

func (rb *record64Bytes) compare(other *record64Bytes) int {
	return bytes.Compare(rb[:], other[:])
}

// B+Tree 中的记录
type record struct {
	keyLen uint8          // offset:0
	key    *record64Bytes // offset:1
	value  *recordUint64  // offset:65
}

type nodeFlag byte

const (
	leafFlag  nodeFlag = 1 << 0
	setupFlag nodeFlag = 1 << 7
)

func (f nodeFlag) isLeaf() bool {
	return f&leafFlag == 1
}

type node struct {
	// flags        nodeFlag byte
	// recordNum    uint8
	// nextPage     GlobalPageID uint64
	// record       [MaxRecordPerPage]record // 最多存储 55 条记录

	nodeID GlobalPageID
	s      space
}

const (
	node1Offset = 0               // nodeFlag
	node2Offset = node1Offset + 1 // recoredNum
	node3Offset = node2Offset + 1 // nextPage
	node4Offset = node3Offset + 8 // record
)

func (t *BPTree) loadNode(id GlobalPageID) *node {
	page, err := t.allocator.Get(id)
	if err != nil {
		return nil
	}
	return &node{
		nodeID: id,
		s:      page.space,
	}
}

func (n *node) dumpRecords() {
	for i := uint64(0); i < n.recordNum(); i++ {
		childKey := n.record(i).key
		childID := n.record(i).value.toUint64()
		fmt.Printf("#[%x]: %d#\t", (*childKey)[:16], childID)
	}
	fmt.Printf("\n")
}

func (n *node) setup(init func(me *node), num uint64, flag nodeFlag) {
	if n.flag()&setupFlag == 1 {
		return
	}
	n.setRecordNum(num)
	n.setFlag(flag | setupFlag)
	init(n)
}

// 传入的新节点为分裂后的右节点
func (n *node) split(newNode *node) (left *node, right *node, leftKey, rightKey *record64Bytes, err error) {
	if n.flag().isLeaf() {
		return n.splitLeaf(newNode)
	}
	return n.splitNonLeaf(newNode)
}

// 分裂叶子节点
// 将当前节点的记录分成两部分，前一部分留在当前节点作为新左节点，后一部分移动到新节点作为新右节点
// 返回新左节点、新右节点、新左节点的最小key、新右节点的最小key
func (n *node) splitLeaf(newNode *node) (left *node, right *node, leftKey, rightKey *record64Bytes, err error) {
	mid := uint64(math.Ceil(float64(n.recordNum()) / 2))

	// new node for right
	newNode.setup(
		func(me *node) {
			n.copyRecords(me, mid, n.recordNum())
		},
		n.recordNum()-mid,
		n.flag(),
	)
	n.clearRecords(mid, n.recordNum())

	rightKey = n.record(mid).key
	right = newNode

	leftKey = n.record(0).key
	left = n

	// 叶子节点要链接下一个叶子节点
	left.setNextPage(right.nodeID)
	return
}

func (n *node) splitNonLeaf(newNode *node) (left *node, right *node, leftKey, rightKey *record64Bytes, err error) {
	mid := uint64(math.Ceil(float64(n.recordNum()) / 2))

	// new node for right
	newNode.setup(
		func(me *node) {
			n.copyRecords(me, mid, n.recordNum())
		},
		n.recordNum()-mid,
		n.flag(),
	)
	n.clearRecords(mid, n.recordNum())

	rightKey = n.record(mid).key
	right = newNode

	leftKey = n.record(0).key
	left = n

	return
}

// shift移动记录
func (n *node) shiftRecord(index uint64, shift int64) {
	if index == n.recordNum() {
		panic(fmt.Sprintf("shiftRecord: index %d out of range", index))
	}
	if shift == 0 {
		return
	}
	if shift < 0 {
		// shift大于index
		if index < uint64(-shift) {
			copy(n.s.buf[node4Offset:], n.s.buf[node4Offset-shift*RecordSize:])
		} else {
			copy(n.s.buf[node4Offset+(index+uint64(shift))*RecordSize:], n.s.buf[node4Offset+index*RecordSize:])
		}
		//清空最后shift个记录
		n.clearRecords(n.recordNum()+uint64(shift), n.recordNum())

	}
	if shift > 0 {
		copy(n.s.buf[node4Offset+(index+uint64(shift))*RecordSize:], n.s.buf[node4Offset+index*RecordSize:])
	}
}

// 清空[start,end)的记录
func (n *node) clearRecords(start uint64, end uint64) {
	if start <= end {
		return
	}
	if end > n.recordNum() {
		end = n.recordNum()
	}
	buf := n.s.buf[node4Offset+start*RecordSize : node4Offset+end*RecordSize]
	clear(buf)
	n.setRecordNum(n.recordNum() - (end - start))
}

// [start,end)的记录追加到dst
func (n *node) copyRecords(dst *node, start uint64, end uint64) {
	if end > n.recordNum() {
		end = n.recordNum()
	}
	copy(dst.s.buf[node4Offset+dst.recordNum()*RecordSize:], n.s.buf[node4Offset+start*RecordSize:node4Offset+end*RecordSize])
}

// 对index处的记录直接赋值
func (n *node) placeRecord(index uint64, rKeyLen uint8, rKey *record64Bytes, rValue uint64) {
	if index >= n.recordNum() {
		return
	}

	copy(n.s.buf[node4Offset+index*RecordSize:], []byte{rKeyLen})
	copy(n.s.buf[node4Offset+index*RecordSize+1:], rKey[:])
	binary.LittleEndian.PutUint64(n.s.buf[node4Offset+index*RecordSize+65:], rValue)
}

func (n *node) appendRecord(rKeyLen uint8, rKey *record64Bytes, rValue uint64) {
	n.setRecordNum(n.recordNum() + 1)
	n.placeRecord(n.recordNum()-1, rKeyLen, rKey, rValue)
}

func (n *node) flag() nodeFlag {
	return nodeFlag(n.s.buf[0])
}

func (n *node) setFlag(flag nodeFlag) {
	n.s.buf[0] = byte(flag)
}

func (n *node) recordNum() uint64 {
	num := uint64(n.s.buf[node2Offset])
	/*
		在实现中，在非叶子节点中，记录数至少为2，因为分裂时，会将需要分裂的节点的 第一个记录 和 中间的记录 移动到父节点，便于节点的指向。
		而判断是否需要分裂的标准是记录数是否大于等于度数，非叶子节点的第一个记录不参与计算。
	*/
	if !n.flag().isLeaf() && num >= 2 {
		return num - 1
	}
	return num
}

func (n *node) setRecordNum(num uint64) {
	n.s.buf[node2Offset] = byte(num)
}

func (n *node) nextPage() GlobalPageID {
	return GlobalPageID(binary.LittleEndian.Uint64(n.s.buf[node3Offset:node4Offset]))
}

func (n *node) setNextPage(nodeID GlobalPageID) {
	binary.LittleEndian.PutUint64(n.s.buf[node3Offset:node4Offset], uint64(nodeID))
}

func (n *node) first() record {
	return n.record(0)
}

func (n *node) last() record {
	return n.record(n.recordNum() - 1)
}

func (n *node) record(index uint64) record {
	if n.recordNum() <= index {
		panic("record index out of range")
	}
	return record{
		keyLen: n.s.buf[node4Offset+index*RecordSize : node4Offset+index*RecordSize+1][0],
		key:    (*record64Bytes)(n.s.buf[node4Offset+index*RecordSize+1 : node4Offset+index*RecordSize+1+64]),
		value:  (*recordUint64)(n.s.buf[node4Offset+index*RecordSize+65 : node4Offset+index*RecordSize+65+8]),
	}
}
