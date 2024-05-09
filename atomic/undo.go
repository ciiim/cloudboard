package atomic

import (
	"io"
	"os"

	"github.com/tinylib/msgp/msgp"
)

var (
	undoStart = func(aID uint64) *UndoLog {
		return &UndoLog{
			AID: aID,
			Op: Op{
				OpID: "start",
			},
		}
	}

	undoEnd = func(aID uint64) *UndoLog {
		return &UndoLog{
			AID: aID,
			Op: Op{
				OpID: "end",
			},
		}
	}
)

type Undo struct {
	// 操作日志路径
	undoLogPath string

	undoLogFile *os.File

	w *msgp.Writer

	r *msgp.Reader

	// 恢复undo日志
	undoGroup *recoverUndoGroup
}

type recoverUndoGroup struct {
	undoMap map[uint64][]Op
}

func NewUndo(undoLogPath string) *Undo {
	return &Undo{
		undoLogPath: undoLogPath,
	}
}

func (u *Undo) OpenUndoLog() error {
	file, err := os.Open(u.undoLogPath)
	if err != nil {
		return err
	}

	u.undoLogFile = file

	u.w = msgp.NewWriter(file)

	return nil
}

func (u *Undo) undoStart(aID uint64) error {
	// 写入开始标志
	return undoStart(aID).EncodeMsg(u.w)
}

func (u *Undo) writeUndo(undo *UndoLog) error {
	// 写undo日志
	return undo.EncodeMsg(u.w)
}

func (u *Undo) undoFinish(aID uint64) error {
	// 写入结束标志
	return undoEnd(aID).EncodeMsg(u.w)
}

func (u *Undo) ClearUndoLog() error {
	return os.Truncate(u.undoLogPath, 0)
}

func (u *Undo) ParseUndoLog() error {
	// 读取undo日志
	if u.r == nil {
		u.r = msgp.NewReader(u.undoLogFile)
	}
	if u.undoGroup == nil {
		u.undoGroup = &recoverUndoGroup{
			undoMap: make(map[uint64][]Op),
		}
	}
	var err error
	for {
		undoLog := &UndoLog{}
		err = undoLog.DecodeMsg(u.r)
		if err != nil {
			break
		}
		if _, ok := u.undoGroup.undoMap[undoLog.AID]; !ok {
			u.undoGroup.undoMap[undoLog.AID] = make([]Op, 0)
		}
		if undoLog.Op.OpID == "end" {
			if undos, ok := u.undoGroup.undoMap[undoLog.AID]; ok {
				if len(undos) > 0 && undos[0].OpID == "start" {
					delete(u.undoGroup.undoMap, undoLog.AID)
				}
			}
		}
		u.undoGroup.undoMap[undoLog.AID] = append(u.undoGroup.undoMap[undoLog.AID], undoLog.Op)
	}
	if err != io.EOF {
		return err
	}
	return u.ClearUndoLog()
}
