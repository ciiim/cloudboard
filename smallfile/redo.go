//go:generate msgp
package smallfile

// var (
// 	ErrTooLarge = errors.New("Data Too Large")
// )

// const (
// 	free  = 0
// 	dirty = 1
// )

// type FreeLog struct {
// 	//span
// 	gPID GlobalPageID
// }

// /*
// redo log
// 防止意外崩溃导致的空闲span未释放
// 和脏页未写回

// [0] free redo

// [1] dirty redo
// */
// type RedoLog struct {
// 	w *msgp.Writer

// 	// 恢复时使用
// 	r *msgp.Reader
// }

// const (
// 	writeIdxStorageSize = 4
// )

// type RedoLogWriter struct {
// 	full []byte
// 	// [0:4] 记录写入位置
// 	idx     []byte
// 	log     []byte
// 	logSize uint32

// 	writeIdx uint32
// }

// type RedoLogReader struct {
// 	log     []byte
// 	logSize uint32

// 	lastWriteIdx uint32

// 	recoverIdx uint32
// }

// func (f *RedoLogWriter) Write(p []byte) (n int, err error) {
// 	if uint32(len(p)) > f.logSize {
// 		return 0, ErrTooLarge
// 	}

// 	if f.writeIdx+f.logSize > f.logSize {
// 		left := uint32(f.logSize) - f.writeIdx
// 		remain := uint32(f.logSize) - left

// 		copy(f.log[f.writeIdx:], p[:left])
// 		copy(f.log[:remain], p[left:remain])

// 		f.writeIdx = uint32(remain)

// 		_ = unix.Msync(f.full, unix.MS_ASYNC)

// 	} else {
// 		copy(f.log[f.writeIdx:], p)

// 		f.writeIdx += f.logSize
// 		if f.writeIdx >= f.logSize {
// 			f.writeIdx = f.writeIdx % f.logSize
// 		}
// 	}
// 	binary.LittleEndian.PutUint32(f.idx, f.writeIdx)
// 	return len(p), nil
// }

// func (f *RedoLogReader) Read(p []byte) (n int, err error) {
// 	if f.lastWriteIdx == f.recoverIdx {
// 		return 0, io.EOF
// 	}

// 	if f.recoverIdx+uint32(len(p)) >= f.lastWriteIdx {
// 		copy(p, f.log[f.recoverIdx:f.lastWriteIdx])
// 		return int(f.lastWriteIdx - f.recoverIdx), io.EOF
// 	}

// 	if f.recoverIdx+uint32(len(p)) > f.logSize {

// 		left := f.logSize - f.recoverIdx
// 		remain := len(p) - int(left)

// 		copy(p, f.log[f.recoverIdx:])
// 		copy(p[left:], f.log[:remain])

// 		f.recoverIdx = uint32(remain)

// 		return len(p), io.EOF
// 	} else {
// 		copy(p, f.log[f.recoverIdx:])
// 		f.recoverIdx += uint32(len(p))
// 		return len(p), nil

// 	}

// }

// func newRedoLog(buf []byte) *RedoLog {
// 	return &RedoLog{
// 		w: msgp.NewWriter(&RedoLogWriter{
// 			full:    buf,
// 			idx:     buf[0:writeIdxStorageSize],
// 			log:     buf[writeIdxStorageSize:],
// 			logSize: uint32(len(buf) - writeIdxStorageSize),
// 		}),
// 	}
// }

// func (f *RedoLog) Recover(b []byte) []FreeLog {
// 	if b[5] == 0 {
// 		return []FreeLog{}
// 	}

// 	lastIdx := binary.LittleEndian.Uint32(b[0:writeIdxStorageSize]) - 1
// 	recoverIdx := func() uint32 {
// 		if lastIdx < uint32(len(b)) && b[lastIdx] == 0 {
// 			return 0
// 		} else if lastIdx+1 >= uint32(len(b)) {
// 			return 0
// 		} else {
// 			return lastIdx + 1
// 		}
// 	}()

// 	f.r = msgp.NewReader(&RedoLogReader{
// 		log:          b[writeIdxStorageSize:],
// 		logSize:      uint32(len(b) - writeIdxStorageSize),
// 		lastWriteIdx: lastIdx,
// 		recoverIdx:   recoverIdx,
// 	})
// 	logs := make([]FreeLog, 0)
// 	var err error
// 	for err == io.EOF {
// 		log := FreeLog{}

// 		err = log.DecodeMsg(f.r)
// 		if err == io.EOF {
// 			logs = append(logs, log)
// 			break
// 		}
// 		if err != nil {
// 			break
// 		}
// 	}
// 	clear(b)
// 	f.r = nil
// 	return logs
// }

// func (f *RedoLog) AddRedo(gPID GlobalPageID, pages int) error {
// 	log := FreeLog{
// 		gPID: gPID,
// 	}

// 	return log.EncodeMsg(f.w)
// }
