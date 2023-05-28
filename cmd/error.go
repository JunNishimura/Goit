package cmd

import "errors"

var (
	ErrGoitNotInitialized = errors.New("not a goit repository (or any of the parent directories): .goit")
	ErrIOHandling         = errors.New("IO handling error")
	ErrInvalidArgs        = errors.New("fatal: invalid arguments")
	ErrIncompatibleFlag   = errors.New("error: incompatible pair of flags")
	ErrNotSpecifiedHash   = errors.New("error: no specified object hash")
	ErrTooManyArgs        = errors.New("error: to many arguments")
	ErrInvalidHash        = errors.New("error: not a valid object hash")
	ErrInvalidHEAD        = errors.New("fatal: could not resolve HEAD")
)
