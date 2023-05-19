package object

import "errors"

var (
	ErrInvalidObject       = errors.New("invalid object")
	ErrInvalidCommitObject = errors.New("invalid commit object")
	ErrNotCommitObject     = errors.New("not commit object")
)
