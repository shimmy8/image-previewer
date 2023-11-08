package imgproxy

import (
	"context"
)

type ImgProxy struct {
}

func New() *ImgProxy {
	//TODO
	return &ImgProxy{}
}

func (iprx *ImgProxy) GetImage(ctx context.Context, url string, headers map[string][]string) ([]byte, error) {
	// TODO
	return nil, nil
}
