package lod

import (
	"errors"

	"golang.org/x/xerrors"
)

var (
	Errf      = xerrors.Errorf
	ErrAs     = errors.As
	ErrIs     = errors.Is
	ErrJoin   = errors.Join
	ErrOpaque = xerrors.Opaque
)
