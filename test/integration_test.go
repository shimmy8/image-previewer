package integration_test

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeRequest(url string) (status int, body []byte) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("%v", err)
	}

	if res.StatusCode != http.StatusOK {
		return res.StatusCode, nil
	}

	resData, _ := io.ReadAll(res.Body)
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
		status, _ := makeRequest("http://localhost:8080/fill/444/555/http:/nginx/500")

		require.Equal(t, http.StatusBadGateway, status)
	})
}
