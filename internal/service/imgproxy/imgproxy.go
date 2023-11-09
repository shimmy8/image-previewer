package imgproxy

import (
	"context"
	"io"
	"net/http"
)

type ImgProxy struct {
}

func New() *ImgProxy {
	return &ImgProxy{}
}

func (iprx *ImgProxy) GetImage(ctx context.Context, url string, headers map[string][]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrResponseNotOk
	}

	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return resData, nil
}
