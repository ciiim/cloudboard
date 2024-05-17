package chash

import (
	"sort"
)

type nilCHashItem struct {
}

func (e *nilCHashItem) ID() string {
	return "nil"
}

func (e *nilCHashItem) Compare(o CHashItem) bool {
	return false
}

type Querier[T CHashItem] struct {
	c *ConsistentHash
}

func NewQuerier[T CHashItem](c *ConsistentHash) *Querier[T] {
	return &Querier[T]{c: c}
}

func (q *Querier[T]) Range(action func(item T)) {
	q.c.hashRingMutex.RLock()
	defer q.c.hashRingMutex.RUnlock()
	for _, v := range q.c.hashMap {
		action(v.(T))
	}
}

func (q *Querier[T]) GetAll(canOutput func(item CHashItem) bool) []T {
	q.c.hashRingMutex.RLock()
	defer q.c.hashRingMutex.RUnlock()
	items := make([]T, 0, q.c.Len())
	for _, v := range q.c.hashMap {
		if canOutput(v) {
			items = append(items, v.(T))
		}
	}
	return items
}

func (q *Querier[T]) GetByID(id string) T {
	q.c.hashRingMutex.RLock()
	defer q.c.hashRingMutex.RUnlock()
	item, ok := q.c.hashMap[q.c.hashFn(string2Bytes(id))]
	if !ok {
		return CHashItem(&nilCHashItem{}).(T)
	}
	return item.(*InnerItem).real.(T)
}

func (q *Querier[T]) Get(key []byte) T {

	q.c.hashRingMutex.RLock()
	defer q.c.hashRingMutex.RUnlock()

	if len(q.c.hashRing) == 0 {
		return CHashItem(&nilCHashItem{}).(T)
	}
	hash := q.c.hashFn(key)

	index := sort.Search(len(q.c.hashRing), func(i int) bool { return q.c.hashRing[i] >= hash })
	item := q.c.hashMap[q.c.hashRing[index%len(q.c.hashRing)]]
	realItem, ok := item.(*InnerItem)
	if ok {
		return realItem.real.(T)
	}
	return CHashItem(&nilCHashItem{}).(T)
}

// 获取key最接近的节点以及后 n 个不相同的节点
// 如果节点数不足 n 个，则返回所有不重复的节点
func (q *Querier[T]) GetN(key []byte, n int) []T {
	if n <= 0 {
		return nil
	}

	q.c.hashRingMutex.RLock()
	defer q.c.hashRingMutex.RUnlock()

	if len(q.c.hashRing) == 0 {
		return nil
	}

	findNMap := make(map[string]CHashItem)

	hash := q.c.hashFn(key)

	// 获取最接近的节点
	index := sort.Search(len(q.c.hashRing), func(i int) bool { return q.c.hashRing[i] >= hash })
	items := make([]T, 0, n)

	items = append(items, q.c.hashMap[q.c.hashRing[index%len(q.c.hashRing)]].(*InnerItem).real.(T))
	findNMap[items[0].ID()] = items[0]

	// 剩余需要遍历的节点数
	remaining := len(q.c.hashRing) - 1

	// 获取后 n-1 个不相同的节点
	index++
	for count := 1; count < n && remaining > 0; func() {
		index++
		remaining--
	}() {
		item := q.c.hashMap[q.c.hashRing[index%len(q.c.hashRing)]].(*InnerItem).real

		// 如果节点已经存在，则跳过
		if compare(item, findNMap) {
			continue
		}

		// 添加节点
		items = append(items, item.(T))
		findNMap[item.ID()] = item
		count++
	}
	return items
}
