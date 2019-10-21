package myasthurts

import "errors"

var (
	ErrTypeNotFound    = errors.New("type not found")
	ErrBuiltInNotFound = errors.New("builtin package not found")
)
