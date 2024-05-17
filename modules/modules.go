package modules

import (
	"errors"
	"sync"

	"github.com/ciiim/cloudborad/modules/modulestype"
)

var (
	ErrModuleNotFound = errors.New("module not found")
)

var (
	Modules *ModuleManager
)

type Module interface {
	Init(params any) error
	Name() string
	Type() modulestype.ModuleType
	Interests() []modulestype.Action
	OnError(err error)
}

/*
MiddlewareModule 中间件模块

拦截文件流，对文件流进行处理，交给下一个中间件。
若没有下一个中间件，则交给最终的功能模块处理。
*/
type MiddlewareModule interface {
	Module
}

/*
FeatureModule 功能模块

不对文件进行直接处理，不拦截文件流。
*/
type FeatureModule interface {
	Module
}

type ModuleManager struct {
	modules sync.Map
}

func new() *ModuleManager {
	return &ModuleManager{
		modules: sync.Map{},
	}
}

func (m *ModuleManager) reg(mod Module) {
	if mod == nil {
		panic("module is nil")
	}
	m.modules.Store(mod.Name(), mod)
}
