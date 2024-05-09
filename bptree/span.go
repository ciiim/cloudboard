package bptree

import (
	"encoding/binary"
	"errors"
)

var (
	ErrOutOfSpace    = errors.New("out of space")
	ErrSpanInUse     = errors.New("span in use")
	ErrNoEnoughPage  = errors.New("no enough page")
	ErrInvaildAccess = errors.New("invaild access")
	ErrInternalError = errors.New("internal error")
)

type spanFlag uint

const (
	spanFreeFlag spanFlag = 0

	spanUsedFlag spanFlag = 1 << 0

	spanHeadFlag spanFlag = 1 << 1
)

func (f spanFlag) isUsed() bool {
	return f&1 == spanUsedFlag
}

func (f spanFlag) isHead() bool {
	return f&spanHeadFlag == spanHeadFlag
}

/*
0-5bit: 5bit flag

6-24bit: 19bits next free span

25-63bit: 39bits span length，实际上不能超过MaxSpanPageNum
*/
type spanInfo uint64

func toSpanInfo(flag spanFlag, spanPages uint, nextFree localPageID) spanInfo {

	if spanPages > MaxSpanPageNum {
		spanPages = MaxSpanPageNum
	}

	// 保证flag只有5bit
	flag &= 0x1F

	return spanInfo((uint(nextFree) << 6) | (spanPages << 25) | uint(flag))
}

func (i spanInfo) nextFree() localPageID {
	return localPageID(i >> 6 & 0x7ffff)
}

func (i spanInfo) spanPages() uint {
	return uint(i >> 25 & 0x7fffff)
}

func (i spanInfo) flag() spanFlag {
	return spanFlag(i & 0x1F)
}

var (
	sysSpan      = toSpanInfo(spanUsedFlag|spanHeadFlag, SysUsedPagePerChunk, SysUsedPagePerChunk)
	initFreeSpan = toSpanInfo(spanHeadFlag, MaxUserPagePerChunk, 0)
)

// 第一个span的nextFree存着最近的一个free span的index
type spanMap struct {
	buf space
}

func (s spanMap) Close() error {
	return s.buf.Release()
}

func (s spanMap) init() {
	s.overwriteSpanHead(0, sysSpan)
	s.overwriteSpanHead(0+SysUsedPagePerChunk, initFreeSpan)
}

func (s spanMap) occupySuitableSpan(pages int) (localPageID, error) {
	//记录上一个span的index
	prev := localPageID(0)
	now := s.getSpanInfo(0).nextFree()
	for {
		span := s.getSpanInfo(now)
		// fmt.Printf("now finding page id:%d,used:%v,head:%v\n",
		// 	now, span.flag().isUsed(), span.flag().isHead())

		// 如果找到的nearest free span已经被使用了，说明chunk已经满了
		if span.flag().isUsed() {
			return 0, ErrNoEnoughPage
		}

		if span.spanPages() >= uint(pages) {
			i := toSpanInfo(spanUsedFlag|spanHeadFlag, uint(pages), 0)
			if err := s.occupySpan(prev, now, i); err != nil {
				return 0, err
			}
			return now, nil
		}

		// 该chunk没有合适的span
		if span.spanPages()+uint(now) >= MaxPagePerChunk {
			return 0, ErrNoEnoughPage
		}

		prev = now
		now = span.nextFree()
	}
}

func (s spanMap) freeSpan(pageID localPageID) (uint, error) {
	if pageID < SysUsedPagePerChunk {
		return 0, ErrInvaildAccess
	}

	span := s.getSpanInfo(pageID)
	if !span.flag().isUsed() {
		return 0, ErrInvaildAccess
	}
	var prev, next localPageID = 0, 0
	//从头开始找到前驱和后继free span
	for {
		next = s.getSpanInfo(prev).nextFree()
		if next == 0 {
			break
		}

		// 找到了
		if next > pageID {
			break
		}

		prev = next
	}

	prevInfo := s.getSpanInfo(prev)
	nextInfo := s.getSpanInfo(next)

	// fmt.Printf("prepare to free pageID:%d,prev:%d,next:%d\n", pageID, prev, next)

	freePages := span.spanPages()

	//链接或合并后继span
	//如果与后继free span相连则合并
	if pageID+localPageID(span.spanPages()) == next {
		// 清空后继span head
		s.overwriteSpanHead(next, toSpanInfo(0, 0, 0))
		span = toSpanInfo(spanFreeFlag|spanHeadFlag, span.spanPages()+nextInfo.spanPages(), nextInfo.nextFree())
		s.overwriteSpanHead(pageID, span)
		// fmt.Printf("merge to next:%d\n", next)

	} else {
		// 否则链接
		// fmt.Printf("link to next:%d\n", next)
		span = toSpanInfo(spanFreeFlag|spanHeadFlag, span.spanPages(), next)
		s.overwriteSpanHead(pageID, span)
	}

	s.tryUpdateNearestFreeSpan(pageID)

	//链接或合并前驱span
	//如果与前驱free span相连则合并
	//不允许和系统span合并
	if prev >= SysUsedPagePerChunk && prev+localPageID(prevInfo.spanPages()) == pageID {
		// 清空当前span head
		s.overwriteSpanHead(pageID, toSpanInfo(0, 0, 0))

		//重写前驱span info
		prevInfo = toSpanInfo(prevInfo.flag(), prevInfo.spanPages()+span.spanPages(), span.nextFree())
		s.overwriteSpanHead(prev, prevInfo)
		// fmt.Printf("merge to prev:%d\n", prev)
	} else {
		// 否则链接
		prevInfo = toSpanInfo(prevInfo.flag(), prevInfo.spanPages(), pageID)
		s.overwriteSpanHead(prev, prevInfo)
		// fmt.Printf("link to prev:%d\n", prev)
	}

	return freePages, nil

}

func (s spanMap) getSpanInfo(pageID localPageID) spanInfo {
	return spanInfo(binary.LittleEndian.Uint64(s.buf.buf[uint(pageID)*uint(s.buf.step) : uint(pageID+1)*uint(s.buf.step)]))
}

func (s spanMap) occupySpan(prevPageID, nowPageID localPageID, newInfo spanInfo) error {

	if (uint(nowPageID)+1)*uint(s.buf.step) > uint(len(s.buf.buf)) {
		return ErrOutOfSpace
	}

	nowFreeSpan := s.getSpanInfo(nowPageID)
	if nowFreeSpan.flag().isUsed() {
		return ErrSpanInUse
	}

	// 只允许从free span的头进行占用
	if !nowFreeSpan.flag().isHead() {
		return ErrInvaildAccess
	}

	freeLength := nowFreeSpan.spanPages()
	if freeLength < newInfo.spanPages() {
		return ErrNoEnoughPage
	}

	// 更新为used span
	s.overwriteSpanHead(nowPageID, newInfo)

	// 生成新的free span
	var newFreeHeadSpan spanInfo
	var newFreeHeadFlag spanFlag = spanHeadFlag
	var newFreeSpanLength = freeLength - newInfo.spanPages()
	var newFreeSpanHeadIndex, newNextFree localPageID

	newFreeSpanHeadIndex = nowPageID + localPageID(newInfo.spanPages())

	// 继承原来的next free
	newNextFree = nowFreeSpan.nextFree()

	if newFreeSpanLength > 0 {
		newFreeHeadSpan = toSpanInfo(newFreeHeadFlag, newFreeSpanLength, newNextFree)
		s.overwriteSpanHead(newFreeSpanHeadIndex, newFreeHeadSpan)

		// prev链上新的free span
		if prevPageID != 0 {

			// fmt.Printf("link to prevPageID:%d\n", prevPageID)

			prevSpan := s.getSpanInfo(prevPageID)
			prevSpan = toSpanInfo(prevSpan.flag(), prevSpan.spanPages(), newFreeSpanHeadIndex)
			s.overwriteSpanHead(prevPageID, prevSpan)
		}

		s.tryUpdateNearestFreeSpan(newFreeSpanHeadIndex)
	} else {
		// 新的free span长度为0，直接指向下一个free span
		// prev链上新的free span
		if prevPageID != 0 {
			prevSpan := s.getSpanInfo(prevPageID)
			prevSpan = toSpanInfo(prevSpan.flag(), prevSpan.spanPages(), newNextFree)
			s.overwriteSpanHead(prevPageID, prevSpan)
		}

		// 如果新的free span不存在，则尝试更新为新的next free span
		s.tryUpdateNearestFreeSpan(newNextFree)
	}

	return nil
}

func (s spanMap) tryUpdateNearestFreeSpan(newNextFree localPageID) {
	span := s.getSpanInfo(0)
	oldNextFree := s.getSpanInfo(span.nextFree())
	if !oldNextFree.flag().isUsed() && newNextFree >= span.nextFree() {
		return
	}
	span = toSpanInfo(span.flag(), span.spanPages(), newNextFree)
	s.overwriteSpanHead(0, span)
	// fmt.Printf("update nearest free span to:%d\n", newNextFree)
}

// 覆写span的头部
func (s spanMap) overwriteSpanHead(pageID localPageID, newInfo spanInfo) {
	binary.LittleEndian.PutUint64(s.buf.buf[uint(pageID)*uint(s.buf.step):uint(pageID+1)*uint(s.buf.step)], uint64(newInfo))
}
