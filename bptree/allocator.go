package bptree

import "fmt"

type Allocator struct {
	d *diskHeap
	c *Cache
}

func NewAllocator() *Allocator {
	d := newDiskHeap()
	writeBack := func(span *SpanInCache) {
		// 写回到磁盘
		if err := d.writeBack(span.globlID, span.buf, span.Pages()); err != nil {
			println("Write back error", err.Error())
		}
	}
	c := NewCache(64, 32, writeBack)
	return &Allocator{d, c}
}

/*
Dump dumps the span info to stdout.

You can set 'skip' to 0 to start the dump from the beginning of the span.
*/
func (a *Allocator) Dump(chunkID chunkID, skip, num int) {
	fmt.Printf("\n--Dump-Chunk-%d--\n", chunkID)
	var span *SpanInCache
	if a.c.busySpanLRU.size > 0 {
		fmt.Printf("  Busy span cache:\n	%-8s %-8s\n", "ID", "Pages")
		for _, v := range a.c.busySpanLRU.m {
			span = v.Value.(*SpanInCache)
			fmt.Printf("	%-8d %-8d\n", span.Id(), span.Pages())
		}
	}
	if a.c.freeSpanLRU.size > 0 {
		fmt.Printf("  Free span cache:\n	%-8s %-8s\n", "ID", "Pages")
		for _, v := range a.c.freeSpanLRU.m {
			span = v.Value.(*SpanInCache)
			fmt.Printf("	%-8d %-8d\n", span.Id(), span.Pages())
		}

	}
	a.d.dump(chunkID, skip+1024, num)
}

func (a *Allocator) AllocNoCache(size int) (*SpanInCache, error) {
	// 从磁盘分配
	id, err := a.d.allocSpan(size)
	if err != nil {
		return nil, err
	}

	return a.d.getSpan(id)
}

func (a *Allocator) Alloc(size int) (GlobalPageID, error) {
	// 从空闲span缓存中分配
	span, err := a.c.allocSpan(size)
	if span != nil && err == nil {
		return span.globlID, nil
	}

	// 从磁盘分配
	return a.d.allocSpan(size)
}

// MarkDirty marks the span as dirty.
// If not, the span will not be written back to disk
func (a *Allocator) MarkDirty(id GlobalPageID) {
	// 标记为脏页
	a.c.markDirty(id)
}

func (a *Allocator) Get(id GlobalPageID) (*SpanInCache, error) {
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

// Force sync the span buffer to disk
func (a *Allocator) ForceSync(id GlobalPageID) {
	// 视作脏页，写回到磁盘
	a.c.busySpanLRU.dirtyWriteBack(id)
}

func (a *Allocator) Free(id GlobalPageID) {
	// 释放span
	if success := a.c.freeSpan(id); success {
		return
	}
	_ = a.d.freeSpan(id)
}

func (a *Allocator) Close() {
	a.c.flush()
	a.d.close()
}

func (a *Allocator) Flush() {
	// 写回所有脏页
	a.c.flush()
}
