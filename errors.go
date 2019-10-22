package myasthurts

import "errors"

var (
	ErrTypeNotFound             = errors.New("type not found")
	ErrBuiltInNotFound          = errors.New("builtin package not found")
	ErrPackageAliasNotFound     = errors.New("package alias not found")
	ErrUnexpectedSelector       = errors.New("unexpected selector identifier")
	ErrUnexpectedExpressionType = errors.New("unexpected expression type")
)
