package bptree

import (
	"errors"
	"sync"
)

type SpanInCache struct {
	globalID GlobalPageID
	//buf
	space space

	dirty bool
}

func (s *SpanInCache) Pages() int {
	return len(s.space.buf) / PageSize
}

func (s *SpanInCache) Id() GlobalPageID {
	return s.globalID
}

func (s *SpanInCache) FixedBytes() []byte {
	return s.space.buf
}

// type SuperCache struct {
// 	partialCache []*Cache
// }

type Cache struct {
	lock sync.Mutex

	// 忙span缓存
	busySpanLRU *lruBusySpanCache

	// 空闲span缓存，可以被分配
	freeSpanLRU *lruFreeSpanCache

	// freeRedo *RedoLog
}

func NewCache(busyCap, freeCap int, writeBackFn func(*SpanInCache), releaseSpanFn func(*SpanInCache)) *Cache {
	c := &Cache{}

	c.busySpanLRU = newBusyLRU(busyCap, writeBackFn)

	c.freeSpanLRU = newFreeLRU(freeCap, releaseSpanFn)
	return c
}

// 标记为脏页
func (c *Cache) markDirty(id GlobalPageID) {
	c.lock.Lock()
	defer c.lock.Unlock()

	span := c.busySpanLRU.getNoUpdate(id)
	if span != nil {
		span.dirty = true
	}
}

func (c *Cache) freeSpan(id GlobalPageID) bool {

	c.lock.Lock()
	defer c.lock.Unlock()

	if span := c.busySpanLRU.getNoUpdate(id); span != nil {
		c.busySpanLRU.remove(id)
		c.freeSpanLRU.put(id, span)
		return true
	}
	return false
}

func (c *Cache) allocSpan(size int) (*SpanInCache, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var pages int

	if size%PageSize != 0 {
		pages = size/PageSize + 1
	} else {
		pages = size / PageSize
	}

	// 从freeSpanLRU中找
	span := c.freeSpanLRU.getClosestSpanSize(pages)
	if span == nil {
		return nil, errors.New("free span cache missing")
	}
	c.freeSpanLRU.remove(span.globalID)
	c.busySpanLRU.put(span.globalID, span)
	return span, nil
}

func (c *Cache) putBusySpan(id GlobalPageID, span *SpanInCache) {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.busySpanLRU.put(id, span)
}

func (c *Cache) putFreeSpan(id GlobalPageID, span *SpanInCache) {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.freeSpanLRU.put(id, span)
}

func (c *Cache) flush() {
	c.busySpanLRU.lock.Lock()
	defer c.busySpanLRU.lock.Unlock()

	for _, v := range c.busySpanLRU.m {
		span := v.Value.(*SpanInCache)
		if span.dirty {
			c.busySpanLRU.writeBack(span)
			span.dirty = false
		}
	}

}

func (c *Cache) getSpan(id GlobalPageID) *SpanInCache {

	c.lock.Lock()
	defer c.lock.Unlock()

	// 从busySpanLRU中找
	span := c.busySpanLRU.get(id)
	return span
}
