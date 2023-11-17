package server

import "errors"

var (
	ErrInvalidRequestParams = errors.New("invalid request params")
	ErrInvalidTargetWidth   = errors.New("invalid target image width")
	ErrInvalidTargetHeight  = errors.New("invalid target image heigth")
	ErrInvalidURL           = errors.New("invalid image URL")
)

// type ResponseErrorCode string

const (
	ErrProxyNotOkResponse     string = "err_get_image_http_not_ok"
	ErrProxyGetImageResponse  string = "err_get_image"
	ErrImageResizeResponse    string = "err_image_resize"
	ErrImageFormatResponse    string = "err_image_format"
	ErrFileNotAnImageResponse string = "err_not_an_image"
	ErrUnknownResponse        string = "err_unknown"
)
