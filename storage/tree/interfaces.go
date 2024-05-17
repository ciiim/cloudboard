package tree

import "github.com/ciiim/cloudborad/storage/types"

type ITreeFileSystem interface {
	AllSpaces() []SpaceInfo
	NewLocalSpace(space string, cap types.Byte) error
	GetLocalSpace(space string) *Space
	DeleteLocalSpace(space string) error
}

// Limited
type Limited interface {
	Occupy() types.Byte
}
