package cmd

import "errors"

var (
	ErrGoitNotInitialized = errors.New("not a goit repository (or any of the parent directories): .goit")
	ErrInvalidArgs        = errors.New("fatal: invalid arguments")
)
