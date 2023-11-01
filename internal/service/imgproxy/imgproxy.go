package imgproxy

import (
	"context"
	"image"
)

type ImgProxy struct {
}

func New() *ImgProxy {
	//TODO
	return &ImgProxy{}
}

func (iprx *ImgProxy) GetImage(ctx context.Context, url string, headers map[string]string) (image.Image, error) {
	// TODO
	return nil, nil
}
