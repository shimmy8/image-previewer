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

func NewHandler(ctx context.Context, app *app.App, logger *zap.Logger) *http.ServeMux {
	h := Handler{app: app, ctx: ctx, logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("/fill/", loggingMiddleware(h.handleResizeRequest, logger))

	return mux
}

func loggingMiddleware(h http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(
			"Reqeust start",
			zap.String("URL", r.URL.String()),
			zap.String("method", r.Method),
			zap.String("User-Agent", r.UserAgent()),
			zap.String("Remote addr", r.RemoteAddr),
		)
		h(w, r)
	}
}

func (h *Handler) handleResizeRequest(w http.ResponseWriter, r *http.Request) {
	// e.g /fill/300/200/url
	pathWithoutRoot := strings.Replace(r.URL.Path, "/fill/", "", 1)
	pathArgs := strings.SplitN(pathWithoutRoot, "/", 3)

	if len(pathArgs) != 3 {
		h.logger.Error("Invalid path args", zap.Strings("parsed", pathArgs), zap.Error(ErrInvalidRequestParams))
		http.Error(w, ErrInvalidRequestParams.Error(), http.StatusBadRequest)
		return
	}

	targetWidth, err := strconv.Atoi(pathArgs[0])
	if err != nil {
		h.logger.Error("Width convertation error", zap.String("parsed width", pathArgs[0]), zap.Error(err))
		http.Error(w, ErrInvalidTargetWidth.Error(), http.StatusBadRequest)
		return
	}

	targetHeight, err := strconv.Atoi(pathArgs[1])
	if err != nil {
		h.logger.Error("Height convertation error", zap.String("parsed height", pathArgs[1]), zap.Error(err))
		http.Error(w, ErrInvalidTargetHeight.Error(), http.StatusBadRequest)
		return
	}

	imgURL, err := url.Parse(pathArgs[2])
	if err != nil {
		h.logger.Error("URL parse error", zap.String("parsed URL", pathArgs[2]), zap.Error(err))
		http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
		return
	}
	if imgURL.Path == "" {
		h.logger.Error("empty URL", zap.String("parsed URL", pathArgs[2]))
		http.Error(w, ErrInvalidURL.Error(), http.StatusBadRequest)
		return
	}

	if imgURL.Scheme == "" {
		imgURL.Scheme = "https"
	}

	result, err := h.app.ResizeImage(h.ctx, r.Header, imgURL.String(), targetWidth, targetHeight)
	if err != nil {
		h.logger.Error(
			"Resize error",
			zap.String("image URL", imgURL.String()),
			zap.Int("width", targetWidth),
			zap.Int("heigth", targetHeight),
			zap.Error(err),
		)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(result); writeErr != nil {
		h.logger.Error("Write response error", zap.Error(writeErr))
	}
}
