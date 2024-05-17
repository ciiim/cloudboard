package replica

import "github.com/ciiim/cloudborad/modules/modulestype"

type ReplicaModule struct {
}

func New() *ReplicaModule {
	return &ReplicaModule{}
}

func (r *ReplicaModule) Name() string {
	return "replica"
}

func (r *ReplicaModule) Type() modulestype.ModuleType {
	return modulestype.Feature
}

func (r *ReplicaModule) Init(params any) error {
	return nil
}

func (r *ReplicaModule) Interests() []modulestype.Action {
	return []modulestype.Action{
		modulestype.Action{
			modulestype.ChunkRead,
			modulestype.ActionError,
		},
		modulestype.Action{},
	}
}
