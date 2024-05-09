package replica

type ReplicaModule struct {
}

func (r *ReplicaModule) Name() string {
	return "replica"
}

func (r *ReplicaModule) Load(params any) error {
	return nil
}
