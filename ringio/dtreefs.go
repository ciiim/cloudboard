package ringio

import (
	"log/slog"

	"github.com/ciiim/cloudborad/node"
	"github.com/ciiim/cloudborad/storage/tree"
	"github.com/ciiim/cloudborad/storage/types"
)

const (
	NEW_SPACE = "__NEW__SPACE__"
)

// Distributed Tree File System.
// Implement FileSystem interface
type DTreeFileSystem struct {
	local  *tree.TreeFileSystem
	remote *rpcTreeClient
	ns     *node.NodeServiceRO
	l      *slog.Logger
}

var _ ITreeDFileSystem = (*DTreeFileSystem)(nil)

func NewDTreeFileSystem(config *tree.Config, ns *node.NodeServiceRO, logger *slog.Logger) *DTreeFileSystem {
	TreeDFileSystem := &DTreeFileSystem{
		local:  tree.NewTreeFileSystem(config),
		remote: newRPCTreeClient(),
		ns:     ns,
		l:      logger,
	}
	return TreeDFileSystem

}

func (dt *DTreeFileSystem) Local() tree.ITreeFileSystem {
	return dt.local
}

func (dt *DTreeFileSystem) pickNode(key []byte) *node.Node {
	return dt.ns.Pick(key)
}

func (dt *DTreeFileSystem) NewSpace(space string, cap types.Byte) error {
	ni := dt.pickNode([]byte(space))
	if dt.ns.Self().Equal(ni) {
		return dt.local.NewLocalSpace(space, cap)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.newSpace(ctx, ni, space, cap)
}

func (dt *DTreeFileSystem) DeleteSpace(space string) error {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		return dt.local.DeleteLocalSpace(space)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.deleteSpace(ctx, ni, space)
}

// 从所有节点收集Space
func (dt *DTreeFileSystem) AllSpaces() []tree.SpaceInfo {
	res := make([]tree.SpaceInfo, 0)
	for _, ni := range dt.ns.GetAllReal() {
		if ni == nil {
			continue
		}
		if ni.Equal(dt.ns.Self()) {
			res = append(res, dt.local.AllSpaces()...)
		} else {
			ctx, cancel := ctxWithTimeout()
			defer cancel()
			spaces, err := dt.remote.allSpaces(ctx, ni)
			if err != nil {
				dt.l.Error(err.Error())
				continue
			}
			res = append(res, spaces...)
		}
	}
	return res
}

func (dt *DTreeFileSystem) GetSpaceStat(space, key string) (*tree.SpaceStatElement, error) {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return nil, tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		space := dt.local.GetLocalSpace(space)
		if space == nil {
			return nil, tree.ErrSpaceNotFound
		}
		return space.GetStatElement(key)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.getSpaceStat(ctx, ni, space, key)
}

func (dt *DTreeFileSystem) SetSpaceStat(space string, e *tree.SpaceStatElement) error {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		space := dt.local.GetLocalSpace(space)
		if space == nil {
			return tree.ErrSpaceNotFound
		}
		return space.SetStatElement(e)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.setSpaceStat(ctx, ni, space, e)

}

func (dt *DTreeFileSystem) MakeDir(space, base, name string) error {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		return dt.local.MakeDir(space, base, name)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.makeDir(ctx, ni, space, base, name)
}

func (dt *DTreeFileSystem) RenameDir(space, base, name, newName string) error {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		if err := dt.local.RenameDir(space, base, name, newName); err != nil {
			return err
		}
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.renameDir(ctx, ni, space, base, name, newName)
}

func (dt *DTreeFileSystem) DeleteDir(space, base, name string) error {
	ni := dt.pickNode([]byte(space))
	if ni == nil {
		return tree.ErrSpaceNotFound
	}
	if dt.ns.Self().Equal(ni) {
		return dt.local.DeleteDir(space, base, name)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.deleteDir(ctx, ni, space, base, name)
}

func (dt *DTreeFileSystem) GetDirSub(space, base, name string) ([]*tree.SubInfo, error) {
	ni := dt.pickNode([]byte(space))

	if dt.ns.Self().Equal(ni) {
		return dt.local.GetDirSub(space, base, name)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.getDirSub(ctx, ni, space, base, name)
}

func (dt *DTreeFileSystem) GetMetadata(space, base, name string) ([]byte, error) {
	ni := dt.pickNode([]byte(space))

	if dt.ns.Self().Equal(ni) {
		return dt.local.GetMetadata(space, base, name)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.getMetadata(ctx, ni, space, base, name)
}

func (dt *DTreeFileSystem) PutMetadata(space, base, name string, fileHash []byte, metadata []byte) error {
	ni := dt.pickNode([]byte(space))

	if dt.ns.Self().Equal(ni) {
		return dt.local.PutMetadata(space, base, name, fileHash, metadata)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.putMetadata(ctx, ni, space, base, name, metadata)
}
func (dt *DTreeFileSystem) DeleteMetadata(space, base, name string) error {
	ni := dt.pickNode([]byte(space))

	if dt.ns.Self().Equal(ni) {
		println("delete metadata local")
		return dt.local.DeleteMetadata(space, base, name)
	}
	ctx, cancel := ctxWithTimeout()
	defer cancel()
	return dt.remote.deleteMetadata(ctx, ni, space, base, name)
}

func (dt *DTreeFileSystem) Node() *node.NodeServiceRO {
	return dt.ns
}