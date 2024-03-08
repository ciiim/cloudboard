package bptree

import (
	"container/list"
	"sync"
)

type lruSpanCache struct {
	size int
	cap  int

	lock sync.Mutex

	l *list.List
	m map[globalPageID]*list.Element
}

type lruBusySpanCache struct {
	*lruSpanCache
	writeBack func(*SpanInCache)
}

type lruFreeSpanCache struct {
	*lruSpanCache
	freeSpan func(globalPageID)
}

func newLRU(cap int) *lruSpanCache {
	return &lruSpanCache{
		size: 0,
		cap:  cap,

		lock: sync.Mutex{},

		l: list.New(),
		m: make(map[globalPageID]*list.Element),
	}
}

func newBusyLRU(cap int, writeBack func(*SpanInCache)) *lruBusySpanCache {
	return &lruBusySpanCache{
		lruSpanCache: newLRU(cap),
		writeBack:    writeBack,
	}
}

func newFreeLRU(cap int, freeSpan func(globalPageID)) *lruFreeSpanCache {
	return &lruFreeSpanCache{
		lruSpanCache: newLRU(cap),
		freeSpan:     freeSpan,
	}
}

func (l *lruBusySpanCache) removeWriteBack(id globalPageID) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.Remove(e)
		delete(l.m, id)
		l.size--
		l.writeBack(e.Value.(*SpanInCache))
	}

}

func (l *lruSpanCache) getClosestSpanSize(pages int) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	for e := l.l.Front(); e != nil; e = e.Next() {
		span := e.Value.(*SpanInCache)
		if span.Pages() >= pages {
			l.l.Remove(e)
			return span
		}
	}
	return nil
}

func (l *lruBusySpanCache) put(id globalPageID, span *SpanInCache) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.MoveToFront(e)
		return nil
	}

	e := l.l.PushFront(span)
	l.m[id] = e
	l.size++

	var ret *SpanInCache = nil

	if l.size > l.cap {
		e := l.l.Back()
		l.l.Remove(e)
		delete(l.m, e.Value.(*SpanInCache).globlID)
		ret = e.Value.(*SpanInCache)
		l.size--
		if ret.dirty {
			l.writeBack(ret)
		}
	}
	return ret

}

func (l *lruSpanCache) put(id globalPageID, span *SpanInCache) {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.MoveToFront(e)
		return
	}

	e := l.l.PushFront(span)
	l.m[id] = e
	l.size++

	if l.size > l.cap {
		e := l.l.Back()
		l.l.Remove(e)
		delete(l.m, e.Value.(*SpanInCache).globlID)
		l.size--
	}
}

func (l *lruSpanCache) getNoUpdate(id globalPageID) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		return e.Value.(*SpanInCache)
	}
	return nil
}

func (l *lruSpanCache) get(id globalPageID) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.MoveToFront(e)
		return e.Value.(*SpanInCache)
	}
	return nil

}

func (l *lruSpanCache) remove(id globalPageID) {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.Remove(e)
		delete(l.m, id)
		l.size--
	}
}
