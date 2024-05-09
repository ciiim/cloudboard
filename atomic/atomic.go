package atomic

import (
	"errors"
	"sync"
	"sync/atomic"
)

type AState int

const (
	AReady AState = iota

	// ADoing 操作进行中
	ADoing

	// AFailed 操作失败
	AFailed

	// ACommitting 操作提交中
	ACommitting

	// ACommitted 操作已提交
	ACommitted

	// ARollbacking 操作回滚中
	ARollbacking

	// ARollbacked 操作已回滚
	ARollbacked
)

const (
	// OpSet 设置操作
	OpSet OpType = iota

	// OpUndo 撤销操作
	OpUndo
)

var (
	toUndo = func(id OpID) OpID {
		return id + "undo"
	}
)

type OpType int8

type OpID string

type OperationFn func(args ...string) error

type operation struct {

	// 操作类型
	OpType OpType

	// 操作id
	OpID OpID

	Args []string

	// 关联函数
	Func OperationFn
}

type AManager struct {
	// 初始化的时候添加操作，后续不可更改，不需要加锁
	operations map[OpID]*operation

	seq atomic.Uint64

	undoManager *Undo
}

func NewAManager() *AManager {
	a := &AManager{
		operations:  make(map[OpID]*operation),
		undoManager: NewUndo("undo.log"),
	}
	_ = a.undoManager.OpenUndoLog()
	return a
}

func (a *AManager) InitOperation(opID OpID, fn OperationFn, undoFn OperationFn) {
	a.operations[opID] = &operation{
		OpType: OpSet,
		OpID:   opID,
		Func:   fn,
	}

	a.operations[toUndo(opID)] = &operation{
		OpType: OpUndo,
		OpID:   opID,
		Func:   undoFn,
	}
}

type Atomic struct {
	seq uint64

	state AState

	opOnce sync.Once
	ops    []operation

	undoOnce sync.Once
	undos    []operation

	manager *AManager
}

func NewAtomic(m *AManager) *Atomic {
	return &Atomic{
		seq: m.seq.Add(1),

		state:   ADoing,
		manager: m,

		ops:   make([]operation, 0),
		undos: make([]operation, 0),
	}
}

/*
重建Atomic对象，从undo日志中恢复

完成后逐个执行对应的undo操作
*/
func (a *AManager) Recover() error {
	if err := a.undoManager.ParseUndoLog(); err != nil {
		return err
	}
	for _, undoLog := range a.undoManager.undoGroup.undoMap {
		for _, op := range undoLog {
			if op.OpID == "start" {
				continue
			}
			fn, ok := a.operations[OpID(op.OpID)]
			if !ok {
				continue
			}
			if err := fn.Func(op.Args...); err != nil {
				return err
			}
		}

	}
	return nil
}

// args[0] 正向操作参数
// args[1] 反向操作参数
func (a *Atomic) RegDo(id OpID, args ...[]string) error {
	// 正向操作
	op, ok := a.manager.operations[id]
	if !ok {
		return errors.New("op not found")
	}
	cop := *op
	cop.Args = args[0]
	a.ops = append(a.ops, cop)

	// 反向操作
	uop, ok := a.manager.operations[toUndo(id)]
	if !ok {
		return errors.New("undo op not found")
	}
	cuop := *uop
	cuop.Args = args[1]
	a.undos = append(a.undos, cuop)

	// 写入undo日志
	a.manager.undoManager.undoStart(a.seq)

	return nil
}

func (a *Atomic) Commit() AState {
	a.opOnce.Do(func() {
		a.state = ACommitting
		for _, op := range a.ops {
			if err := op.Func(op.Args...); err != nil {
				a.state = AFailed
				return
			}
		}
		a.state = ACommitted

		// 删除undo日志
		a.manager.undoManager.undoFinish(a.seq)
	})
	return a.state
}

func (a *Atomic) RollBack() AState {
	a.undoOnce.Do(func() {
		a.state = ARollbacking
		for _, op := range a.undos {
			if err := op.Func(op.Args...); err != nil {
				a.state = AFailed
			}
		}
		a.state = ARollbacked

		// 删除undo日志
		a.manager.undoManager.undoFinish(a.seq)
	})
	return a.state
}
