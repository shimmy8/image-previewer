package resizer

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	"golang.org/x/image/bmp"
)

type Resizer struct{}

func New() *Resizer {
	return &Resizer{}
}

func (r *Resizer) ResizeImage(
	imageBytes []byte,
	targetWitdh int,
	targetHeigth int,
) ([]byte, error) {
	img, imgFmt, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			return nil, fmt.Errorf("%w: %w", ErrNotAnImage, err)
		}
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
	if encodeErr := encoder(buf, res); encodeErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Resizer) encodeJPEG(w io.Writer, m image.Image) error {
	return jpeg.Encode(w, m, nil)
}
