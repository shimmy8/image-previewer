package server

import "errors"

var (
	ErrInvalidRequestParams = errors.New("invalid request params")
	ErrInvalidTargetWidth   = errors.New("invalid target image width")
	ErrInvalidTargetHeight  = errors.New("invalid target image heigth")
	ErrInvalidURL           = errors.New("invalid image URL")
)
