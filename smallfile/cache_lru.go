package smallfile

import (
	"container/list"
	"sync"
)

type lruSpanCache struct {
	size int
	cap  int

	lock sync.Mutex

	l *list.List
	m map[GlobalPageID]*list.Element
}

type lruBusySpanCache struct {
	*lruSpanCache
	writeBack func(*SpanInCache)
}

type lruFreeSpanCache struct {
	*lruSpanCache
	releaseSpan func(*SpanInCache)
}

func newLRU(cap int) *lruSpanCache {
	return &lruSpanCache{
		size: 0,
		cap:  cap,

		lock: sync.Mutex{},

		l: list.New(),
		m: make(map[GlobalPageID]*list.Element),
	}
}

func newBusyLRU(cap int, writeBack func(*SpanInCache)) *lruBusySpanCache {
	return &lruBusySpanCache{
		lruSpanCache: newLRU(cap),
		writeBack:    writeBack,
	}
}

func newFreeLRU(cap int, releaseSpan func(*SpanInCache)) *lruFreeSpanCache {
	return &lruFreeSpanCache{
		lruSpanCache: newLRU(cap),
		releaseSpan:  releaseSpan,
	}
}

func (l *lruBusySpanCache) dirtyWriteBack(id GlobalPageID) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.writeBack(e.Value.(*SpanInCache))
	}

}

func (l *lruFreeSpanCache) getClosestSpanSize(pages int) *SpanInCache {

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

func (l *lruBusySpanCache) put(id GlobalPageID, span *SpanInCache) {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.MoveToFront(e)
		return
	}

	e := l.l.PushFront(span)
	l.m[id] = e
	l.size++

	var ret *SpanInCache = nil

	if l.size > l.cap {
		e := l.l.Back()
		l.l.Remove(e)
		delete(l.m, e.Value.(*SpanInCache).globalID)
		ret = e.Value.(*SpanInCache)
		l.size--
		if ret.dirty {
			l.writeBack(ret)
		}
		_ = ret.space.Release()
	}
}

func (l *lruFreeSpanCache) put(id GlobalPageID, span *SpanInCache) *SpanInCache {

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
		delete(l.m, e.Value.(*SpanInCache).globalID)
		ret = e.Value.(*SpanInCache)
		l.size--
		l.releaseSpan(ret)
	}
	return ret

}

func (l *lruSpanCache) getNoUpdate(id GlobalPageID) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		return e.Value.(*SpanInCache)
	}
	return nil
}

func (l *lruSpanCache) get(id GlobalPageID) *SpanInCache {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.MoveToFront(e)
		return e.Value.(*SpanInCache)
	}
	return nil

}

func (l *lruSpanCache) remove(id GlobalPageID) {

	l.lock.Lock()
	defer l.lock.Unlock()

	if e, ok := l.m[id]; ok {
		l.l.Remove(e)
		delete(l.m, id)
		l.size--
	}
}
