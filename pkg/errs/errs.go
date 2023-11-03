package errs

import "errors"

var (
	ErrInternal      = errors.New("internal error")
	ErrRender        = errors.New("failed to render")
	ErrInvalidObject = errors.New("invalid object")
)
