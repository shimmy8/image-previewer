package resizer

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
)

type Resizer struct {
}

func New() *Resizer {
	return &Resizer{}
}

func (r *Resizer) ResizeImage(ctx context.Context, imageBytes []byte, targetWitdh int, targetHeigth int) ([]byte, error) {
	img, imgFmt, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	var encoder func(w io.Writer, m image.Image) error

	switch imgFmt {
	case "png":
		encoder = png.Encode
	case "jpeg":
		encoder = r.encodeJPEG
	case "bmp":
		encoder = bmp.Encode
	default:
		return nil, ErrFormatNotSupported
	}

	res := imaging.Fill(img, targetWitdh, targetHeigth, imaging.Center, imaging.Lanczos)

	buf := new(bytes.Buffer)
	encodeErr := encoder(buf, res)
	if encodeErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Resizer) encodeJPEG(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}
