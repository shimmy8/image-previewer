package integration_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func makeRequest(url string) (status int, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
	}

	resData, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	return res.StatusCode, resData
}

func loadLocalImage(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	bytes := make([]byte, fileInfo.Size())

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return bytes
}

func TestSuccessResize(t *testing.T) {
	t.Run("Test success downsize", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/333/666/http:/nginx/images/_gopher_original_1024x504.jpg")
		expected := loadLocalImage("./testdata/images/gopher_333x666.jpg")

		require.Equal(t, http.StatusOK, status)
		require.Equal(t, expected, body)
	})

	t.Run("Test success upsize", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/images/gopher_50x50.jpg")
		expected := loadLocalImage("./testdata/images/gopher_444x555_shakal.jpg")

		require.Equal(t, http.StatusOK, status)
		require.Equal(t, expected, body)
	})
}

func TestResizeErrors(t *testing.T) {
	t.Run("test image not found", func(t *testing.T) {
		status, _ := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/images/gopher_800x50.jpg")

		require.Equal(t, http.StatusBadGateway, status)
	})

	t.Run("test image server error", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/500")

		require.Equal(t, http.StatusBadGateway, status)
		require.Equal(t, []byte("err_get_image_http_not_ok\n"), body)
	})

	t.Run("test server unavailable", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/444/555/not-available-server.com/image.png")

		require.Equal(t, http.StatusBadGateway, status)
		require.Equal(t, []byte("err_get_image\n"), body)
	})

	t.Run("test image format not supported", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/images/monkey_90x90.gif")

		require.Equal(t, http.StatusBadGateway, status)
		require.Equal(t, []byte("err_image_format\n"), body)
	})

	t.Run("test not an image", func(t *testing.T) {
		status, body := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/not_image")

		require.Equal(t, http.StatusBadGateway, status)
		require.Equal(t, []byte("err_not_an_image\n"), body)
	})

	t.Run("test cached image", func(t *testing.T) {
		sourceImage := loadLocalImage("./testdata/images/gopher_500x500.png")
		imgCopyFilename := "gopher_copy.png"
		imgCopyFullName := filepath.Join("./testdata/images/", imgCopyFilename)

		// create copy of a file at the folder, available to nginx
		writeErr := os.WriteFile(imgCopyFullName, sourceImage, os.ModePerm)
		defer func() {
			// cleanup if smth wnt wrong during test an file was not deleted
			if _, err := os.Stat(imgCopyFullName); errors.Is(err, os.ErrNotExist) {
				return
			}
			os.Remove(imgCopyFullName)
		}()

		require.NoError(t, writeErr)

		nginxFileURL := fmt.Sprintf("http://localhost:80/images/%s", imgCopyFilename)
		// make sure file created
		ngnixStatus, _ := makeRequest(nginxFileURL)
		require.Equal(t, http.StatusOK, ngnixStatus)

		resizeURL := fmt.Sprintf(
			"http://localhost:8080/fill/100/50/http:/nginx/images/%s",
			imgCopyFilename,
		)

		// request file size change with new filename
		status, _ := makeRequest(resizeURL)
		require.Equal(t, http.StatusOK, status)

		// now delete file from disk
		os.Remove(imgCopyFullName)
		// make sure file deleted
		ngnixRepStatus, _ := makeRequest(nginxFileURL)
		require.Equal(t, http.StatusNotFound, ngnixRepStatus)

		// requiest a resize again
		repStatus, _ := makeRequest(resizeURL)
		// voila! cached file is still there
		require.Equal(t, http.StatusOK, repStatus)
	})
}
