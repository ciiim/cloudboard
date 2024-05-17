package smallfile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	stack "github.com/duke-git/lancet/v2/datastructure/stack"
)

const (
	// UserPageNo 是分配器分配的第一个页面，存储的是根节点的GlobalPageID
	NodeDataPageID = UserPageNo
)

const (
	IndexFileName = "index.db"
)

var (
	ErrNeedSplit    = errors.New("bptree: need split")
	ErrCannotMerge  = errors.New("bptree: cannot merge")
	ErrCannotBorrow = errors.New("bptree: cannot borrow")
	ErrNotFound     = errors.New("bptree: not found")
)

/*
B+Tree

非叶子节点的value存储的是下一个节点的GlobalPageID

下一节点的key都大于等于当前节点的key

只有叶子节点的value存储的是数据的GlobalPageID

索引节点指向的数据节点的key都大于等于当前节点的key
*/
type BPTree struct {
	allocator *Allocator

	//最大度数，最多存储degree-1个key，当节点内的记录数达到或大于degree时，需要分裂。
	degree int

	root GlobalPageID
}

func NewBPTree(degree int) *BPTree {
	b := &BPTree{
		allocator: NewAllocator(IndexFileName),
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
	id, err := t.allocator.Alloc(PageSize)
	if err != nil {
		panic("b+ tree init error:" + err.Error())
	}
	if id != NodeDataPageID {
		panic("b+ tree init error: id not match:" + fmt.Sprintf("returns id:%d,expected %d", id, UserPageNo))
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

func (t *BPTree) Dump(depth int) {
	if t.root == 0 {
		fmt.Printf("root is nil\n")
		return
	}
	rootNode := t.loadNode(t.root)
	fmt.Printf("root id:%d\troot record num:%d\tleaf:%s\n", t.root, rootNode.recordNum(), func() string {
		if rootNode.flag().isLeaf() {
			return "true"
		}
		return "false"
	}())
	rootNode.dumpRecords()
}

func (t *BPTree) Insert(key []byte, dataPageID GlobalPageID) error {
	return t.insertRecord(record{
		keyLen: uint8(len(key)),
		key:    (*record64Bytes)(key),
		value:  GlobalPageIDToBytes(dataPageID),
	})
}

func (t *BPTree) Delete(key []byte) error {
	return t.deleteRecord(key)
}

func (t *BPTree) Search(key []byte) (GlobalPageID, bool) {
	if t.root == 0 {
		println("root is 0")
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
		node = t.loadNode(GlobalPageID(nextID.toUint64()))
	}
	index, found := t.binarySearch(keyp, node)
	if !found {
		return 0, false
	}
	return GlobalPageID(node.record(index).value.toUint64()), true
}

func (t *BPTree) GetData(gPID GlobalPageID) (*dataSpan, error) {
	return nil, nil
}

/*
make root 创建根
分裂原有根
*/
func (t *BPTree) makeRoot() (*node, error) {
	newRootID, err := t.allocator.Alloc(PageSize)
	if err != nil {
		return nil, err
	}

	// 如果是空树
	if t.root == 0 {
		t.root = newRootID
		newRoot := t.loadNode(newRootID)
		// 置为叶子节点
		newRoot.setFlag(leafFlag)
		return newRoot, t.saveRoot()
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

	// 置为叶子节点
	left.setFlag(leafFlag)
	right.setFlag(leafFlag)

	newRoot := t.loadNode(newRootID)
	// 更新根节点
	newRoot.placeRecord(0, 0, leftKey, uint64(left.nodeID))
	newRoot.placeRecord(1, 0, rightKey, uint64(right.nodeID))

	// 置为非叶子节点
	newRoot.setFlag(newRoot.flag() & ^leafFlag)

	t.root = newRootID

	return newRoot, t.saveRoot()
}

type treeStackElement struct {
	nodeID      GlobalPageID
	recordIndex uint64
}

/*
插入记录

1.查找插入位置
2.插入记录
如果节点记录数大于degree，分裂节点
如果插入位置为第一个，递归更新父节点的key
*/
func (t *BPTree) insertRecord(record record) error {

	// 空树
	if t.root == 0 {
		newRoot, err := t.makeRoot()
		if err != nil {
			return err
		}
		_, err = t.insertRecordToNode(newRoot, record)
		return err
	}

	// 查找栈
	// 用于节点分裂时更新父节点
	stack := stack.NewArrayStack[treeStackElement]()

	// 从根节点开始查找插入位置
	node := t.loadNode(t.root)
	if node == nil {
		return fmt.Errorf("insertRecord: load node %d failed", t.root)
	}

	// 查找至叶子节点
	for !node.flag().isLeaf() {
		index, found := t.binarySearch(record.key, node)
		if !found {
			if index == 0 {
				return ErrNotFound
			} else {
				index = index - 1
			}
		}
		stack.Push(treeStackElement{
			nodeID:      node.nodeID,
			recordIndex: index,
		})
		nextID := node.record(index).value.toUint64()
		node = t.loadNode(GlobalPageID(nextID))
	}

	// 插入到叶子节点
	index, err := t.insertRecordToNode(node, record)
	if err != nil && err != ErrNeedSplit {
		return err
	}

	var needUpdateParentKey, needSplit bool = false, false
	// 如果插入位置为第一个，递归更新父节点的key
	if index == 0 {
		needUpdateParentKey = true
	}
	// 如果节点记录数大于degree，分裂节点
	if err == ErrNeedSplit {
		needSplit = true
	}

	// 分裂
	var parentID GlobalPageID
	var parentNodeInfo *treeStackElement
	for needSplit || needUpdateParentKey {
		parentNodeInfo, _ = stack.Pop()
		// 只有stack为空一种情况
		if parentNodeInfo == nil {
			parentID = 0 //不存在的页号，让loadNode返回nil
		} else {
			parentID = parentNodeInfo.nodeID
		}
		parent := t.loadNode(parentID)
		//分裂
		if needSplit {
			fmt.Printf("split node:%d\n", node.nodeID)
			if err = t.splitNode(node, parent); err != ErrNeedSplit && err != nil {
				return err
			} else if err == ErrNeedSplit {
				needSplit = true
			} else {
				needSplit = false
			}
		}
		//更新父节点的key
		if needUpdateParentKey {
			if parent == nil {
				needUpdateParentKey = false
				continue
			}
			parent.placeRecord(parentNodeInfo.recordIndex, record.keyLen, record.key, node.nodeID.toUint64())
			if parentNodeInfo.recordIndex == 0 {
				needUpdateParentKey = true
			} else {
				needUpdateParentKey = false
			}
		}
		node = parent
	}
	fmt.Printf("insert to node:%d\n", node.nodeID)
	return nil
}

func (t *BPTree) insertRecordToNode(node *node, record record) (uint64, error) {
	index, err := t.insertRecordToNodeDirectA(node, record)
	if err != nil {
		return 0, err
	}
	if node.recordNum() >= uint64(t.degree) {
		return index, ErrNeedSplit
	}
	return index, nil
}

// parent 传入nil表示无父节点
func (t *BPTree) splitNode(target, parent *node) error {
	if parent == nil {
		fmt.Printf("split root\n")
		if _, err := t.makeRoot(); err != nil {
			return err
		}
		return nil
	}

	newNode, err := t.allocator.Alloc(PageSize)
	if err != nil {
		return err
	}
	_, right, _, rightKey, err := target.split(t.loadNode(newNode))
	if err != nil {
		return err
	}

	// 更新父节点
	if _, err = t.insertRecordToNodeDirectB(parent, 0, rightKey, right.nodeID.toUint64()); err != nil {
		return err
	}

	//如果父节点满了，继续分裂
	if parent.recordNum() >= uint64(t.degree) {
		return ErrNeedSplit
	}
	return nil
}

// 插入记录到指定节点

func (t *BPTree) insertRecordToNodeDirectA(node *node, record record) (uint64, error) {
	index, found := t.binarySearch(record.key, node)
	if found {
		return 0, nil
	}
	node.setRecordNum(node.recordNum() + 1)
	// 原记录移位
	node.shiftRecord(index, 1)

	//插入记录
	node.placeRecord(index, record.keyLen, record.key, binary.LittleEndian.Uint64(record.value[:]))

	return index, nil
}

func (t *BPTree) insertRecordToNodeDirectB(node *node, rKeyLen uint8, rKey *record64Bytes, rValue uint64) (uint64, error) {
	index, found := t.binarySearch(rKey, node)
	if found {
		return 0, nil
	}
	node.setRecordNum(node.recordNum() + 1)
	// 原记录移位
	node.shiftRecord(index, 1)

	//插入记录
	node.placeRecord(index, rKeyLen, rKey, rValue)

	return index, nil
}

/*
删除记录

1.查找记录
2.删除记录
如果节点记录数小于ceil(M/2)，合并节点
将需要合并的节点中的记录加入到附近的兄弟节点
阶为3的B+树，degree=4
若左节点记录数少于ceil(M/2)，将左节点记录加入到右节点，

如果兄弟节点记录数大于ceil(M/2)，借一个记录

若删除的记录为第一个记录，需要递归更新父节点的key
*/
func (t *BPTree) deleteRecord(key []byte) error {
	// 查找栈
	// 用于节点合并时更新父节点
	treeStack := stack.NewArrayStack[treeStackElement]()

	// 从根节点开始查找插入位置
	node := t.loadNode(t.root)
	if node == nil {
		return fmt.Errorf("deleteRecord: load node %d failed", t.root)
	}

	// 查找至叶子节点
	for !node.flag().isLeaf() {
		index, found := t.binarySearch((*record64Bytes)(key), node)
		if !found {
			if index == 0 {
				return errors.New("deleteRecord: binarySearch error")
			} else {
				index = index - 1
			}
		}
		treeStack.Push(treeStackElement{
			nodeID:      node.nodeID,
			recordIndex: index,
		})
		nextID := node.record(index).value.toUint64()
		node = t.loadNode(GlobalPageID(nextID))
	}

	// 删除叶子节点
	index, err := t.deleteRecordFromNode(node, key)
	if err != nil {
		return err
	}

	var needUpdateParentKey, needMergeOrBorrow bool = false, false
	// 如果删除的记录为第一个，递归更新父节点的key
	if index == 0 {
		needUpdateParentKey = true
	}

	// 如果节点记录数小于ceil(M/2)，合并节点或借记录
	if node.recordNum() <= uint64(math.Ceil(float64(t.degree-1)/2)) {
		needMergeOrBorrow = true
	}

	/*
		四种情况 1.合并 2.借记录 3.只更新父节点的key 4.无事发生
	*/
	if needMergeOrBorrow {
		//尝试合并
		if err = t.mergeNode(treeStack, node); err != nil && err != ErrCannotMerge {
			return err
		} else if err == nil {
			//合并成功，返回
			return nil
		}
		//尝试借记录
		if err = t.borrowRecord(treeStack, node); err != nil {
			return err
		}
	} else if needUpdateParentKey {
		var parentID GlobalPageID
		var parentNodeInfo *treeStackElement
		for needUpdateParentKey {
			parentNodeInfo, _ = treeStack.Pop()
			// 只有stack为空一种情况
			if parentNodeInfo == nil {
				parentID = 0 //不存在的页号，让loadNode返回nil
			} else {
				parentID = parentNodeInfo.nodeID
			}
			parent := t.loadNode(parentID)
			//更新父节点的key
			if parent == nil {
				needUpdateParentKey = false
				continue
			}
			parent.placeRecord(parentNodeInfo.recordIndex, node.first().keyLen, node.first().key, node.nodeID.toUint64())
			if parentNodeInfo.recordIndex == 0 {
				needUpdateParentKey = true
			} else {
				needUpdateParentKey = false
			}
			node = parent
		}
	}
	return nil
}

/*
合并节点

检查兄弟节点，先检查左边再检查右边。

如果左边有节点且记录数小于M，则将待合并节点中的记录加入到左边节点，删除parent中的待合并节点对应的key-value对

如果右边有节点且记录数小于M，则将待合并节点中的记录加入到右边节点，删除parent中的待合并节点对应的key-value对，同时更新合并后的节点的parent中对应的key-value对

不可能没有兄弟节点，因为父节点的记录数不会小于ceil(M/2)

FIXME: 如果合并后父节点不满足B+树特性且为根节点，合并后的节点成为新的根节点
*/
func (t *BPTree) mergeNode(treeStack *stack.ArrayStack[treeStackElement], node *node) error {
	//不满足合并条件直接返回
	if node.recordNum() >= uint64(math.Ceil(float64(t.degree-1)/2)) {
		return ErrCannotMerge
	}
	var parentID GlobalPageID
	var parentNodeInfo *treeStackElement
	parentNodeInfo, _ = treeStack.Pop()
	// 只有stack为空一种情况
	if parentNodeInfo == nil {
		parentID = 0 //不存在的页号，让loadNode返回nil
	} else {
		parentID = parentNodeInfo.nodeID
	}
	parent := t.loadNode(parentID)
	recordIndexInParent := parentNodeInfo.recordIndex
	indexNeedUpdate := recordIndexInParent
	targetNode := node

	// 检查左边
	if recordIndexInParent > 0 {
		leftIndex := recordIndexInParent - 1
		leftID := GlobalPageID(parent.record(leftIndex).value.toUint64())
		left := t.loadNode(leftID)
		if left.recordNum()+node.recordNum() < uint64(t.degree-1) {
			// 将待合并节点中的记录追加到左节点
			node.copyRecords(left, 0, node.recordNum())
			// 删除parent中的待合并节点对应的key-value对
			parent.shiftRecord(leftIndex, -1)

			//释放待合并节点
			t.allocator.Free(node.nodeID)
		}
	} else if recordIndexInParent+1 <= parent.recordNum() {
		rightIndex := recordIndexInParent + 1
		rightID := GlobalPageID(parent.record(rightIndex).value.toUint64())
		right := t.loadNode(rightID)
		if right.recordNum()+node.recordNum() < uint64(t.degree-1) {
			// 右节点腾出位置给待合并节点
			right.shiftRecord(0, int64(node.recordNum()))
			// 将待合并节点中的记录加入到右边节点
			node.copyRecords(right, 0, node.recordNum())
			// 删除parent中的待合并节点对应的key-value对
			parent.shiftRecord(rightIndex, -1)

			targetNode = right
			indexNeedUpdate = rightIndex

			//释放待合并节点
			t.allocator.Free(node.nodeID)
		}
	} else {
		// 把pop出来的放回去
		treeStack.Push(*parentNodeInfo)
		return ErrCannotMerge
	}

	//更新父节点的key
	for indexNeedUpdate == 0 && parent != nil {
		parent.placeRecord(indexNeedUpdate, targetNode.first().keyLen, targetNode.first().key, targetNode.nodeID.toUint64())
		// 父节点上的key不在第一个，不需要递归更新
		if indexNeedUpdate != 0 {
			break
		}
		targetNode = parent
		parentNodeInfo, _ = treeStack.Pop()
		// 只有stack为空一种情况
		if parentNodeInfo == nil {
			break
		}
		parentID = parentNodeInfo.nodeID
		parent = t.loadNode(parentID)
		indexNeedUpdate = parentNodeInfo.recordIndex
	}
	return nil
}

/*
借记录

当删除操作破坏了B+树的性质时，需要借记录

从左边借或者从右边借

如果从左节点借到了记录，一定是位于本节点的第一个，需要递归更新本节点的父节点对应的key

如果从右节点借到了记录，需要递归更新右节点的父节点对应的key

无论如何，都要更新父节点的key
*/
func (t *BPTree) borrowRecord(treeStack *stack.ArrayStack[treeStackElement], node *node) error {
	var parentID GlobalPageID
	var parentNodeInfo *treeStackElement
	var needUpdate bool
	parentNodeInfo, _ = treeStack.Pop()
	// 只有stack为空一种情况
	if parentNodeInfo == nil {
		parentID = 0 //不存在的页号，让loadNode返回nil
	} else {
		parentID = parentNodeInfo.nodeID
	}
	parent := t.loadNode(parentID)
	recordIndexInParent := parentNodeInfo.recordIndex
	indexNeedUpdate := recordIndexInParent
	targetNode := node
	// 检查左边
	if recordIndexInParent > 0 {
		leftIndex := recordIndexInParent - 1
		leftID := GlobalPageID(parent.record(leftIndex).value.toUint64())
		left := t.loadNode(leftID)
		if left.recordNum() > uint64(math.Ceil(float64(t.degree-1)/2)) {
			// 借左节点的最后一个记录
			node.shiftRecord(0, 1)
			borrowedRecord := left.last()
			node.placeRecord(0, borrowedRecord.keyLen, borrowedRecord.key, borrowedRecord.value.toUint64())

			// 删除左边节点的记录
			left.shiftRecord(left.recordNum()-1, -1)

			needUpdate = true
		}
		//检查右边
	} else if recordIndexInParent+1 != parent.recordNum() {
		rightIndex := recordIndexInParent + 1
		rightID := GlobalPageID(parent.record(rightIndex).value.toUint64())
		right := t.loadNode(rightID)
		if right.recordNum() > uint64(math.Ceil(float64(t.degree-1)/2)) {
			// 借一个记录
			borrowedRecord := right.first()
			node.appendRecord(borrowedRecord.keyLen, borrowedRecord.key, borrowedRecord.value.toUint64())

			// 删除右节点的记录
			right.shiftRecord(0, -1)

			// 指向右边节点
			indexNeedUpdate++
			// 目标更新节点指向右节点
			targetNode = right

			needUpdate = true
		}
	} else {
		// 把pop出来的放回去
		treeStack.Push(*parentNodeInfo)
		return ErrCannotBorrow
	}

	//更新父节点的key
	for needUpdate && parent != nil {
		parent.placeRecord(indexNeedUpdate, targetNode.first().keyLen, targetNode.first().key, targetNode.nodeID.toUint64())
		// 父节点上的key不在第一个，不需要递归更新
		if indexNeedUpdate != 0 {
			break
		}
		targetNode = parent
		parentNodeInfo, _ = treeStack.Pop()
		// 只有stack为空一种情况
		if parentNodeInfo == nil {
			break
		}
		parentID = parentNodeInfo.nodeID
		parent = t.loadNode(parentID)
		indexNeedUpdate = parentNodeInfo.recordIndex
	}

	return nil
}

func (t *BPTree) deleteRecordFromNode(node *node, key []byte) (uint64, error) {
	index, found := t.binarySearch((*record64Bytes)(key), node)
	if !found {
		return 0, nil
	}

	// 删除记录
	node.shiftRecord(index, -1)

	return index, nil
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
