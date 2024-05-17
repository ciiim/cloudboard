package smallfile

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/sys/unix"
)

type diskHeap struct {
	file *os.File

	fileSize atomic.Int64

	findChunk    sync.Mutex
	currentChunk chunkID

	spanInfoSpaceMap map[chunkID]spanMap

	metadata *headerMetadata
	busyPage *headerBusyPage
}

func newDiskHeap(fileName string) *diskHeap {
	d := &diskHeap{
		metadata:         new(headerMetadata),
		busyPage:         new(headerBusyPage),
		spanInfoSpaceMap: make(map[chunkID]spanMap),
	}
	if err := d.init(fileName); err != nil {
		return nil
	}
	return d
}

func (d *diskHeap) dump(chunkID chunkID, skip, num int) {
	spanMap, err := d.getSpanMap(chunkID)
	if err != nil {
		fmt.Printf("dump chunk %d failed:%v\n", chunkID, err)
	}
	if d.fileSize.Load() < d.chunkOffset(chunkID)+SpanMapLengthPerChunk {
		fmt.Printf("chunk %d not exist\n", chunkID)
		return
	}
	spanInfo := spanMap.getSpanInfo(localPageID(0))
	fmt.Printf("--span head info--\n  next free:%d\n", spanInfo.nextFree())

	num += skip
	for i := skip; i < num; i++ {
		spanInfo := spanMap.getSpanInfo(localPageID(i))
		fmt.Printf("--page [%d] info--\n  next free:%d\n  span pages:%d\n  span head:%v\n  span used:%v\n",
			i, spanInfo.nextFree(), spanInfo.spanPages(), spanInfo.flag().isHead(), spanInfo.flag().isUsed())
	}
}

/*
WARNING: 任何Chunk Page操作前都要使用growChunkIfNeeded，防止越界访问触发abort
*/
func (d *diskHeap) growChunkIfNeeded(chunkID chunkID, pageID localPageID, pages int) error {
	// fmt.Printf("grow chunk %d, page %d, pages %d\n", chunkID, pageID, pages)
	offset := d.chunkOffset(chunkID) + int64(pageID+localPageID(pages))*PageSize

	return d.growTo(offset)
}

func (d *diskHeap) init(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	d.file = file

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()

	d.fileSize.Store(size)

	// 第一次使用，需要初始化
	if d.fileSize.Load() < ChunkSpaceOffset {
		err = d.grow(ChunkSpaceOffset)
		if err != nil {
			return err
		}
	}

	d.metadata.space, err = d.mapSpace(0, HeaderMetadataByteNum, syscall.MADV_NORMAL, 8)
	if err != nil {
		return err
	}

	d.busyPage.buf, err = d.mapSpace(BusyPageOffset, BusyPageByteNum, syscall.MADV_NORMAL, 4)
	if err != nil {
		return err
	}

	runtime.SetFinalizer(d, func(d *diskHeap) {
		for _, buf := range d.spanInfoSpaceMap {
			buf.Close()
		}
		d.file.Close()
	})
	return nil
}

func (d *diskHeap) close() error {
	runtime.SetFinalizer(d, nil)
	for _, buf := range d.spanInfoSpaceMap {
		buf.Close()
	}
	if err := d.metadata.space.Release(); err != nil {
		return err
	}
	if err := d.busyPage.buf.Release(); err != nil {
		return err
	}
	return d.file.Close()
}

func (d *diskHeap) grow(length int64) error {

	if err := syscall.Fallocate(int(d.file.Fd()), 0, d.fileSize.Load(), length); err != nil {
		return err
	}
	d.fileSize.Add(length)

	// println("grow", length)
	return nil
}

func (d *diskHeap) growTo(offset int64) error {
	if d.fileSize.Load() < offset {
		return d.grow(offset - d.fileSize.Load())
	}
	return nil
}

func (d *diskHeap) mapSpace(offset int64, length int64, advice int, step ...uint) (space, error) {

	var s uint

	if len(step) == 0 {
		s = 8
	} else {
		s = step[0]
	}

	buf, err := syscall.Mmap(int(d.file.Fd()), offset, int(length), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return space{}, err
	}

	_ = syscall.Madvise(buf, advice)

	return space{buf: buf, step: s}, nil
}

func (d *diskHeap) initChunk(id chunkID) error {
	// spanMap
	spanMap, err := d.getSpanMap(id)
	if err != nil {
		return err
	}

	// grow span map
	if err = d.growChunkIfNeeded(id, 0, SysUsedPagePerChunk); err != nil {
		return err
	}

	//init span map
	spanMap.init()

	// busyPage
	d.busyPage.setChunkBusyPages(id, SysUsedPagePerChunk)

	d.metadata.SetNextChunk(id + 1)

	return nil

}

func (d *diskHeap) allocSpanInChunk(chunk chunkID, pages int) (GlobalPageID, error) {
	//读取该chunk的span map
	sm, err := d.getSpanMap(chunk)
	if err != nil {
		return 0, err
	}

	// find d suitable free span
	pageID, err := sm.occupySuitableSpan(pages)
	if err != nil {
		return 0, err
	}

	//update busy page num
	d.busyPage.addChunkBusyPages(chunk, int32(pages))

	gPID := pageID.toGlobal(chunk)

	if err := d.growChunkIfNeeded(chunk, pageID, pages); err != nil {
		return 0, err
	}

	return gPID, nil
}

func (d *diskHeap) chunkOffset(id chunkID) int64 {
	return ChunkSpaceOffset + int64(id)*MaxChunkSize
}

func (d *diskHeap) allocSpan(size int) (gPID GlobalPageID, err error) {
	if size == 0 {
		return 0, ErrInvaildParam
	}

	pageNum := func() int {
		i := size % PageSize
		if i != 0 {
			return int(size/PageSize) + 1
		}
		return int(size / PageSize)
	}()

	// fmt.Printf("try allocate %d page(s)\n", pageNum)

	if pageNum > MaxSpanPageNum {
		return 0, ErrMaxPagePerSpan
	}

	d.findChunk.Lock()
	defer d.findChunk.Unlock()
	var used int
	var looped bool = false
	for {
		used = int(d.busyPage.getChunk(d.currentChunk))
		if used == -1 {
			// busy chunk metadata is full
			return 0, ErrMaxChunkNum
		}
		//used == 0 说明遇到了没有初始化的chunk，如果遍历过一次，说明前面的chunk都没有符合的，需要新初始化一个chunk
		if used == 0 {
			if looped {
				// chunk初始化
				_ = d.initChunk(d.metadata.GetNextChunk())
				used = int(d.busyPage.getChunk(d.currentChunk))
			} else {
				looped = true
				d.currentChunk = 0
				continue
			}
		}
		if MaxPagePerChunk-used < pageNum {
			// chunk 没有足够的空间
			d.currentChunk++
			continue
		}
		gPID, err = d.allocSpanInChunk(d.currentChunk, pageNum)
		return
	}
}

// 获取chunk 的 span info map
func (d *diskHeap) getSpanMap(chunkID chunkID) (spanMap, error) {
	if sm, ok := d.spanInfoSpaceMap[chunkID]; ok {
		return sm, nil
	}

	space, err := d.mapSpace(d.chunkOffset(chunkID), SpanMapLengthPerChunk, syscall.MADV_NORMAL)
	if err != nil {
		return spanMap{}, err
	}

	sm := spanMap{space}

	d.spanInfoSpaceMap[chunkID] = sm

	return sm, nil
}

// 使用msync强制写回
func (d *diskHeap) writeBack(space space) error {
	return unix.Msync(space.buf, unix.MS_ASYNC)
}

func (d *diskHeap) freeSpan(gPID GlobalPageID) error {
	chunkID := gPID.chunkID()
	localID := gPID.toLocal()

	// println("disk free span", gPID)

	var sm spanMap

	//读取该chunk的span map
	sm, err := d.getSpanMap(chunkID)
	if err != nil {
		return err
	}
	freePages := uint(0)
	if freePages, err = sm.freeSpan(localID); err != nil {
		return err
	}

	d.busyPage.addChunkBusyPages(chunkID, -int32(freePages))

	return nil

}

// func (d *diskHeap) getSpanInfo(chunkID chunkID) (spanMap, error) {
// 	if spanMap, ok := d.spanInfoSpaceMap[chunkID]; ok {
// 		return spanMap, nil
// 	}

// 	space, err := d.mapSpace(d.chunkOffset(chunkID), SpanMapLengthPerChunk, syscall.MADV_NORMAL)
// 	if err != nil {
// 		return spanMap{}, err
// 	}

// 	spanInfoMap := spanMap{space}

// 	d.spanInfoSpaceMap[chunkID] = spanInfoMap

// 	return spanInfoMap, nil
// }

func (d *diskHeap) getSpan(gPID GlobalPageID) (*SpanInCache, error) {
	chunkID := gPID.chunkID()

	// 访问不存在的span
	if chunkID >= d.metadata.GetNextChunk() {
		return nil, ErrInvaildAccess
	}

	spanMap, err := d.getSpanMap(gPID.chunkID())
	if err != nil {
		return nil, err
	}
	// 获取信息
	info := spanMap.getSpanInfo(gPID.toLocal())

	if gPID.toLocal() < SysUsedPagePerChunk {
		return nil, ErrInvaildAccess
	}

	// 访问不存在或未分配的span
	if info == 0 || !info.flag().isUsed() {
		return nil, ErrInvaildAccess
	}

	// mmap 需要的span
	spanSpace, err := d.mapSpace(d.chunkOffset(chunkID)+int64(gPID.toLocal())*PageSize, int64(info.spanPages())*PageSize, syscall.MADV_NORMAL)
	if err != nil {
		return nil, err
	}

	span := &SpanInCache{
		globalID: gPID,

		space: spanSpace,

		dirty: false,
	}

	return span, nil
}
