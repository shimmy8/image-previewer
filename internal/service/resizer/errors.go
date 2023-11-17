package resizer

import "errors"

var (
	ErrFormatNotSupported = errors.New("format not supported")
	ErrNotAnImage         = errors.New("not an image")
)
