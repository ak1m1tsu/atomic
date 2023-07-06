package storage

import (
	"errors"
)

var (
	ErrAliasNotFound   = errors.New("alias not found")
	ErrAliasExists     = errors.New("alias already exists")
	ErrInvalidArgument = errors.New("*sql.DB is nil")
)
