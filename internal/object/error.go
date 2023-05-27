package object

import "errors"

var (
	ErrInvalidObject       = errors.New("invalid object")
	ErrInvalidTreeObject   = errors.New("invalid tree object")
	ErrInvalidCommitObject = errors.New("invalid commit object")
	ErrNotCommitObject     = errors.New("not commit object")
	ErrIOHandling          = errors.New("IO handling error")
)
