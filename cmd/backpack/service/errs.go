package service

import "errors"

// Space Error
var (
	ErrSpaceExist      = errors.New("ErrSpaceExist")
	ErrSpaceNotExisted = errors.New("ErrSpaceNotExisted")
	ErrPassword        = errors.New("ErrPassword")
	ErrParam           = errors.New("ErrParam")
)
