package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestServerErrors(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockApp := app.New(10, "./", 2, logger)

	mockHandler := Handler{app: mockApp, logger: logger}

	t.Run("test request params errors", func(t *testing.T) {
		tests := []struct {
			URL        string
			StatusCode int
			Error      error
		}{
			{URL: "/fill/", StatusCode: 400, Error: ErrInvalidRequestParams},
			{URL: "/fill/200/", StatusCode: 400, Error: ErrInvalidRequestParams},
			{URL: "/fill/200/300", StatusCode: 400, Error: ErrInvalidRequestParams},
			{URL: "/fill/20a/300/test", StatusCode: 400, Error: ErrInvalidTargetWidth},
			{URL: "/fill/200/30i/test", StatusCode: 400, Error: ErrInvalidTargetHeight},
			{URL: "/fill/200/300/", StatusCode: 400, Error: ErrInvalidURL},
		}

		for _, tt := range tests {
			tt := tt
			t.Run("request to "+tt.URL, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tt.URL, nil)
				rec := httptest.NewRecorder()

				mockHandler.handleResizeRequest(rec, req)

				res := rec.Result()
				defer res.Body.Close()

				require.Equal(t, tt.StatusCode, res.StatusCode)
				require.Equal(t, tt.Error.Error(), strings.TrimRight(rec.Body.String(), "\n"))
			})
		}
	})
}
