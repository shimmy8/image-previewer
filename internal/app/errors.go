package app

import "errors"

var (
	ErrProxyResponseNotOk = errors.New("image proxy server response not ok")
	ErrProxyGetImage      = errors.New("image proxy error")
	ErrImageFormat        = errors.New("image format not supported")
	ErrFileNotAnImage     = errors.New("file is not an image")
	ErrResizeImage        = errors.New("resize image error")
)
