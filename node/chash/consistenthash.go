package chash

import (
	"hash/crc64"
	"slices"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

type ConsistentHashFn func([]byte) uint64

type CHashItem interface {
	ID() string
	Compare(o CHashItem) bool
}

// 用于实现向后查找真实节点
type InnerItem struct {

	//virtual 是否是虚拟节点
	virtual bool

	//real 真实节点
	real CHashItem
}

func warp2InnerItem(real CHashItem, virtual bool) *InnerItem {
	return &InnerItem{
		virtual: virtual,
		real:    real,
	}
}

func (i *InnerItem) ID() string {
	return i.real.ID()
}

func (i *InnerItem) Compare(item CHashItem) bool {
	return i.real.ID() == item.ID()
}

func (i *InnerItem) IsVirtual() bool {
	return i.virtual
}

// Consistent hash Map
type ConsistentHash struct {
	count int64

	// hashFn hash函数
	hashFn ConsistentHashFn

	//replicas 虚拟节点个数
	replicas int

	hashRingMutex sync.RWMutex

	//hash ring 包含虚拟节点
	hashRing []uint64

	//node info map 包含虚拟节点
	hashMap map[uint64]CHashItem
}

var (
	DefaultHashFn = func(b []byte) uint64 {
		return crc64.Checksum(b, crc64.MakeTable(crc64.ISO))
	}
)

func NewConsistentHash(replicas int, fn ConsistentHashFn) *ConsistentHash {
	m := &ConsistentHash{
		replicas:      replicas,
		hashFn:        fn,
		hashRingMutex: sync.RWMutex{},
		hashRing:      make([]uint64, 0),
		hashMap:       make(map[uint64]CHashItem),
	}

	if fn == nil {
		m.hashFn = DefaultHashFn
	}
	return m
}

// []byte READ ONLY
func string2Bytes(s string) (readOnly []byte) {
	sd := unsafe.StringData(s)
	return unsafe.Slice(sd, len(s))
}

func (c *ConsistentHash) Len() int64 {
	return atomic.LoadInt64(&c.count)
}

func (c *ConsistentHash) Add(item CHashItem) {
	c.hashRingMutex.Lock()
	defer c.hashRingMutex.Unlock()

	//添加真实节点
	hashid := c.hashFn(string2Bytes(item.ID()))

	//添加真实节点到hashMap
	c.hashMap[hashid] = warp2InnerItem(item, false)

	//添加真实节点到hashRing
	c.hashRing = append(c.hashRing, hashid)

	//添加虚拟节点
	for i := 0; i < c.replicas; i++ {
		hashid := c.hashFn(string2Bytes(strconv.Itoa(i) + item.ID()))
		c.hashMap[hashid] = warp2InnerItem(item, true)
		c.hashRing = append(c.hashRing, hashid)
	}
	slices.Sort[[]uint64](c.hashRing)

	atomic.AddInt64(&c.count, 1)
}

func (c *ConsistentHash) Del(item CHashItem) {
	c.hashRingMutex.Lock()
	defer c.hashRingMutex.Unlock()

	//删除真实节点
	hash := c.hashFn(string2Bytes(item.ID()))
	i, found := sort.Find(len(c.hashRing), func(i int) int {
		return func() int {
			if c.hashRing[i] < hash {
				return -1
			} else if c.hashRing[i] == hash {
				return 0
			} else {
				return 1
			}
		}()
	})
	if !found {
		return
	}
	c.hashRing = append(c.hashRing[:i], c.hashRing[i+1:]...)
	delete(c.hashMap, hash)

	for i := 0; i < c.replicas; i++ {
		hash := c.hashFn(string2Bytes(strconv.Itoa(i) + item.ID()))
		delete(c.hashMap, hash)
		//删除虚拟节点
		i, found := sort.Find(len(c.hashRing), func(i int) int {
			return func() int {
				if c.hashRing[i] < hash {
					return -1
				} else if c.hashRing[i] == hash {
					return 0
				} else {
					return 1
				}
			}()
		})
		if !found {
			continue
		}
		c.hashRing = append(c.hashRing[:i], c.hashRing[i+1:]...)
	}

	atomic.AddInt64(&c.count, -1)
}

func compare(a CHashItem, m map[string]CHashItem) bool {
	_, ok := m[a.ID()]
	return ok
}
