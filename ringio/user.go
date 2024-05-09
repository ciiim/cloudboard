package ringio

import "github.com/spf13/afero"

var _ afero.Fs = (*UserSpace)(nil)

type UserSpace struct {
	spaceName string

	spaceBasePath string
}

func (u *UserSpace) Open(name string) (afero.File, error) {

}
