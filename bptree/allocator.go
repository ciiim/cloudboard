package bptree

type Allocator struct {
	d *diskHeap
	c *Cache
}

func NewAllocator() *Allocator {
	d := newDiskHeap()
	writeBack := func(span *SpanInCache) {
		// 写回到磁盘
		_ = d.writeBack(span.globlID, span.buf, span.Pages())
	}
	c := NewCache(64, 16, writeBack)
	return &Allocator{d, c}
}

func (a *Allocator) Dump(chunkID chunkID, skip, num int) {
	a.d.dump(chunkID, skip, num)
}

func (a *Allocator) AllocNoCache(size int) (*SpanInCache, error) {
	// 从磁盘分配
	id, err := a.d.allocSpan(size)
	if err != nil {
		return nil, err
	}

	return a.d.getSpan(id)
}

func (a *Allocator) Alloc(size int) (globalPageID, error) {
	// 从空闲span缓存中分配
	span, err := a.c.allocSpan(size)
	if span != nil && err == nil {
		return span.globlID, nil
	}

	// 从磁盘分配
	return a.d.allocSpan(size)
}

func (a *Allocator) MarkDirty(id globalPageID) {
	// 标记为脏页
	a.c.markDirty(id)
}

func (a *Allocator) Get(id globalPageID) (*SpanInCache, error) {
	// 从繁忙span缓存中找
	span := a.c.getSpan(id)
	if span != nil {
		return span, nil
	}
	// 从磁盘找
	span, err := a.d.getSpan(id)
	if err != nil {
		return nil, err
	}

	// 放入繁忙span缓存
	a.c.putBusySpan(id, span)

	return span, nil
}

func (a *Allocator) ForceSync(id globalPageID) {
	// 强制释放span
	a.c.busySpanLRU.removeWriteBack(id)
}

func (a *Allocator) Free(id globalPageID) {
	// 释放span
	if success := a.c.freeSpan(id); !success {
		_ = a.d.freeSpan(id)
	}
}

func (a *Allocator) Flush() {
	// 写回所有脏页
	a.c.flush()
}
