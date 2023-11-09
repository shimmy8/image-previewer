package resizer

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func loadImage(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return bytes
}

func TestResizeImages(t *testing.T) {
	resizer := New()

	t.Run("test resize jpeg", func(t *testing.T) {
		jpegImgBytes := loadImage("../../../testdata/images/gopher_50x50.jpg")
		jpegSmallImgBytes := loadImage("../../../testdata/images/gopher_20x20.jpg")

		res, err := resizer.ResizeImage(context.Background(), jpegImgBytes, 20, 20)

		require.NoError(t, err)
		require.Equal(t, jpegSmallImgBytes, res)
	})

	t.Run("test resize png", func(t *testing.T) {
		pngImgBytes := loadImage("../../../testdata/images/gopher_500x500.png")
		pngSmallImgBytes := loadImage("../../../testdata/images/gopher_100x100.png")

		res, err := resizer.ResizeImage(context.Background(), pngImgBytes, 100, 100)

		require.NoError(t, err)
		require.Equal(t, pngSmallImgBytes, res)
	})

	t.Run("test resize bmp", func(t *testing.T) {
		bmpImgBytes := loadImage("../../../testdata/images/gopher_320x240.bmp")
		bmpSmallImgBytes := loadImage("../../../testdata/images/gopher_100x100.bmp")

		res, err := resizer.ResizeImage(context.Background(), bmpImgBytes, 100, 100)

		require.NoError(t, err)
		require.Equal(t, bmpSmallImgBytes, res)
	})

	t.Run("test err not supported", func(t *testing.T) {
		gifImgBytes := loadImage("../../../testdata/images/monkey_90x90.gif")

		_, err := resizer.ResizeImage(context.Background(), gifImgBytes, 10, 10)

		require.ErrorIs(t, err, ErrFormatNotSupported)
	})

	t.Run("test err not image", func(t *testing.T) {
		gifImgBytes := loadImage("../../../testdata/not_image.txt")

		_, err := resizer.ResizeImage(context.Background(), gifImgBytes, 10, 10)

		require.ErrorIs(t, err, image.ErrFormat)
	})
}
