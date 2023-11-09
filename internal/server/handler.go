package server

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/shimmy8/image-previewer/internal/app"
	"go.uber.org/zap"
)

type Handler struct {
	app    *app.App
	ctx    context.Context
	logger *zap.Logger
}

func NewHandler(app *app.App, ctx context.Context, logger *zap.Logger) *http.ServeMux {
	h := Handler{app: app, ctx: ctx, logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/fill/", loggingMiddleware(h.handleResizeRequest, logger))

	return mux
}

func loggingMiddleware(h http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("New request", zap.String("URL", r.URL.String()), zap.String("method", r.Method))
		h(w, r)
	}
}

func (h *Handler) handleResizeRequest(w http.ResponseWriter, r *http.Request) {
	// e.g /fill/300/200/url
	pathWithoutRoot := strings.Replace(r.URL.Path, "/fill/", "", 1)
	pathArgs := strings.SplitN(pathWithoutRoot, "/", 3)

	h.logger.Info("path args", zap.Strings("args", pathArgs))
	if len(pathArgs) != 3 {
		http.Error(w, ErrInvalidRequestParams.Error(), http.StatusBadRequest)
		return
	}

	targetWidth, err := strconv.Atoi(pathArgs[0])
	if err != nil {
		http.Error(w, ErrInvalidTargetWidth.Error(), http.StatusBadRequest)
		return
	}

	targetHeight, err := strconv.Atoi(pathArgs[1])
	if err != nil {
		http.Error(w, ErrInvalidTargetHeight.Error(), http.StatusBadRequest)
		return
	}

	imgUrl, err := url.Parse(pathArgs[2])
	if err != nil {
		http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
		return
	}
	if imgUrl.Host == "" || imgUrl.Path == "" {
		http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
		return
	}

	if imgUrl.Scheme == "" {
		imgUrl.Scheme = "https"
	}

	result, err := h.app.ResizeImage(h.ctx, r.Header, imgUrl.String(), targetWidth, targetHeight)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Write(result)
}
