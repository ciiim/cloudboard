package modulestype

type ModuleType int

/*
ActionType 操作类型

模块处理的操作类型，用于区分不同的操作。

模块对感兴趣的操作类型进行监听，当有对应的操作发生时，模块会被调用。
*/
type ActionType int
type ActionSubType int

type Action struct {
	Type    ActionType
	SubType ActionSubType
}

const (
	ChunkRead ActionType = iota
	ChunkNew
	ChunkCountAdd
	ChunkCountMinus
	ChunkDelete

	FileRead
	FileNew
	FileDelete

	NodeJoin
	NodeAlive
	NodeDead
)

const (
	ActionEnter ActionSubType = iota
	ActionDone
	ActionSuccess
	ActionError
)

const (
	Middleware ModuleType = iota
	Feature
)
