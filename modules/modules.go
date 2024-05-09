package modules

import (
	"errors"
	"sync"
)

var (
	ErrModuleNotFound = errors.New("module not found")
)

var (
	Modules *ModuleManager
)

type ActionType int
type ModuleType int

const (
	ChunkRead ActionType = iota
	ChunkWrite
	ChunkDelete

	FileRead
	FileWrite
	FileDelete

	NodeJoin
	NodeAlive
	NodeDead
)

const (
	Middleware ModuleType = iota
	Feature
)

type Module interface {
	Load(params any) error
	Name() string
	OnError(err error)
}

// MiddlewareModule 中间件模块
type MiddlewareModule interface {
	Module
}

// FeatureModule 功能模块
type FeatureModule interface {
}

type ModuleManager struct {
	modules sync.Map
}

func new() *ModuleManager {
	return &ModuleManager{
		modules: sync.Map{},
	}
}

func init() {
	Modules = new()
	Modules.reg()
}

func (m *ModuleManager) reg(mod Module) {
	if mod == nil {
		panic("module is nil")
	}
	m.modules.Store(mod.Name(), mod)
}
