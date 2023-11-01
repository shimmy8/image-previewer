package resizer

import (
	"context"
	"image"
)

type Resizer struct {
}

func New() *Resizer {
	//TODO
	return &Resizer{}
}

func (r *Resizer) ResizeImage(ctx context.Context, image image.Image, targetWitdh int, targetHeigth int) (image.Image, error) {
	//TODO
	return nil, nil
}
