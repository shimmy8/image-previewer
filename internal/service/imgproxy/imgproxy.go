package imgproxy

import (
	"context"
	"io"
	"net/http"
	"time"
)

type ImgProxy struct {
	timeout int
}

func New(timeout int) *ImgProxy {
	return &ImgProxy{timeout: timeout}
}

func (iprx *ImgProxy) GetImage(ctx context.Context, url string, headers map[string][]string) ([]byte, error) {
	rqCtx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(iprx.timeout))
	defer cancel()

	req, err := http.NewRequestWithContext(rqCtx, http.MethodGet, url, nil)
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
