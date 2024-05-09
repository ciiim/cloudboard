package bptree

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	stack "github.com/duke-git/lancet/v2/datastructure/stack"
)

const (
	// UserPageNo 是分配器分配的第一个页面，存储的是根节点的GlobalPageID
	NodeDataPageID = UserPageNo
)

var (
	ErrNeedSplit = errors.New("bptree: need split")
)

type data struct {
	gPID   GlobalPageID
	length int64 // 数据长度
}

type dataSpan struct {
	limitReader *io.LimitedReader
	flags       nodeFlag
}

/*
B+Tree

非叶子节点的value存储的是下一个节点的GlobalPageID

下一节点的key都大于等于当前节点的key

只有叶子节点的value存储的是数据的GlobalPageID

索引节点指向的数据节点的key都大于等于当前节点的key
*/
type BPTree struct {
	allocator *Allocator

	degree int //alias max record num per node + 1
	root   GlobalPageID
}

func NewBPTree(degree int) *BPTree {
	b := &BPTree{
		allocator: NewAllocator(),
		degree:    degree,
	}

	b.load()

	return b
}

func (t *BPTree) load() {
	if page, err := t.allocator.Get(NodeDataPageID); err != nil {
		if errors.Is(err, ErrInvaildAccess) {
			t.initRootDataPage()
		}
	} else {
		t.root = t.loadRoot(page.space.buf)
	}
}

func (t *BPTree) initRootDataPage() {
	id, err := t.allocator.Alloc(1)
	if err != nil {
		panic("b+ tree init error:" + err.Error())
	}
	if id != NodeDataPageID {
		panic("b+ tree init error: id not match")
	}
}

func (t *BPTree) loadRoot(page []byte) GlobalPageID {
	return GlobalPageID(binary.LittleEndian.Uint64(page[:8]))
}

func (t *BPTree) saveRoot() error {
	page, err := t.allocator.Get(NodeDataPageID)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint64(page.space.buf[:8], uint64(t.root))
	t.allocator.ForceSync(NodeDataPageID)
	return nil
}

func (t *BPTree) InsertReader(key []byte, r io.Reader) error {
	return nil
}

func (t *BPTree) InsertData(key []byte, data []byte) error {
	return nil
}

func (t *BPTree) Delete(key []byte) error {
	return nil
}

func (t *BPTree) Search(key []byte) (uint64, bool) {
	if t.root == 0 {
		return 0, false
	}
	keyp := (*record64Bytes)(key)
	// 从root开始
	node := t.loadNode(t.root)
	for !node.flag().isLeaf() {
		index, found := t.binarySearch(keyp, node)
		if !found {
			if index == 0 {
				return 0, false
			} else {
				index = index - 1
			}
		}
		nextID := node.record(index).value
		node = t.loadNode(GlobalPageID(binary.LittleEndian.Uint64(nextID[:8])))
	}
	index, found := t.binarySearch(keyp, node)
	if !found {
		return 0, false
	}
	return binary.LittleEndian.Uint64(node.record(index).value[:8]), true
}

func (t *BPTree) GetData(gPID GlobalPageID) (*dataSpan, error) {
	return nil, nil
}

/*
make root 创建根
*/
func (t *BPTree) makeRoot() (*node, error) {
	newRootID, err := t.allocator.Alloc(PageSize)
	if err != nil {
		return nil, err
	}

	// 如果是空树
	if t.root == 0 {
		t.root = newRootID
		return t.loadNode(newRootID), t.saveRoot()
	}

	oldRoot := t.loadNode(t.root)

	rightID, err := t.allocator.Alloc(PageSize)
	if err != nil {
		return nil, err
	}

	left, right, leftKey, rightKey, err := oldRoot.split(t.loadNode(rightID))
	if err != nil {
		return nil, err
	}

	newRoot := t.loadNode(newRootID)
	// 更新根节点
	newRoot.placeRecord(0, 0, leftKey, uint64(left.nodeID))
	newRoot.placeRecord(1, 0, rightKey, uint64(right.nodeID))

	t.root = newRootID

	return newRoot, t.saveRoot()
}

// 插入记录
func (t *BPTree) insertRecord(record record) error {

	// 空树
	if t.root == 0 {
		newRoot, err := t.makeRoot()
		if err != nil {
			return err
		}
		return t.insertRecordToNode(newRoot, record)
	}

	// 查找栈
	// 用于节点分裂时更新父节点
	stack := stack.NewArrayStack[GlobalPageID]()

	// 从根节点开始查找插入位置
	node := t.loadNode(t.root)
	if node == nil {
		return fmt.Errorf("insertRecord: load node %d failed", node.nodeID)
	}
	stack.Push(node.nodeID)

	// 查找至叶子节点
	for !node.flag().isLeaf() {
		index, found := t.binarySearch(record.key, node)
		if !found {
			if index == 0 {
				return errors.New("insertRecord: binarySearch error")
			} else {
				index = index - 1
			}
		}
		nextID := node.record(index).value.toUint64()
		node = t.loadNode(GlobalPageID(nextID))
		stack.Push(node.nodeID)
	}

	// 插入到叶子节点
	for insertError := t.tryInsertRecordToNode(node, record); insertError == ErrNeedSplit; {
		// 分裂节点

		// 申请节点
		newNodeID, err := t.allocator.Alloc(PageSize)
		if err != nil {
			return err
		}
		newNode := t.loadNode(newNodeID)
		left, right, leftKey, rightKey, err := node.split(newNode)
		if err != nil {
			return err
		}
		parentNodeID, err := stack.Pop()
		/*
			函数内部只有一种错误情况，即栈为空
			if s.IsEmpty() {
				return nil, errors.New("stack is empty")
			}
		*/
		if err != nil {
			// 直接插入后分裂
			t.insertRecordToNode(t.loadNode(*parentNodeID))
			_, err := t.makeRoot()
			if err != nil {
				return err
			}
			break
		}

	}
	return nil
}

// 尝试插入记录到指定节点
// 若节点已满，返回ErrNeedSplit
func (t *BPTree) tryInsertRecordToNode(node *node, record record) error {
	if node.recordNum() >= uint64(t.degree-1) {
		return ErrNeedSplit
	}
	return t.insertRecordToNode(node, record)
}

// 插入记录到指定节点
func (t *BPTree) insertRecordToNode(node *node, record record) error {
	index, found := t.binarySearch(record.key, node)
	if found {
		return nil
	}

	// 原记录移位
	node.shiftRecord(index, 1)

	//插入记录
	node.placeRecord(index, record.keyLen, record.key, binary.LittleEndian.Uint64(record.value[:]))
	return nil
}

// 插入：如果found == false, 那么index是插入的位置，要把index和index之后的所有记录都往后移动一位
//
// 查找：如果found == false，index == 0 说明key比所有记录都小
// 如果found == false，index != 0 说明key在index-1和index之间
// 如果found == false，index == node.recordNum() 说明key比所有记录都大
func (t *BPTree) binarySearch(key *record64Bytes, node *node) (index uint64, found bool) {
	left, right := uint64(0), node.recordNum()

	for left < right {
		index = (right + left) / 2
		midKey := node.record(index).key
		cmp := key.compare(midKey)
		if cmp == 0 {
			return index, true
		}
		if cmp > 0 {
			left = index + 1
		} else {
			right = index
		}
	}
	return left, false
}
