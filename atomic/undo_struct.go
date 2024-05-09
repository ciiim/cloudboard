//go:generate msgp
package atomic

// 一个undo操作
// 记录在undo日志中
type UndoLog struct {
	// 原子操作ID
	AID uint64

	// 操作
	Op Op
}

// 一个操作
type Op struct {

	// 操作id
	// 用于恢复时查找对应的操作函数
	OpID string //OpID
	// 操作参数
	Args []string
}
