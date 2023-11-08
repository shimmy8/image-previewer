package resizer

import (
	"context"
)

type Resizer struct {
}

func New() *Resizer {
	//TODO
	return &Resizer{}
}

func (r *Resizer) ResizeImage(ctx context.Context, imageBytes []byte, targetWitdh int, targetHeigth int) ([]byte, error) {
	//TODO
	return nil, nil
}
