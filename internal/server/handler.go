package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/shimmy8/image-previewer/internal/app"
	"github.com/shimmy8/image-previewer/internal/service/imgproxy"
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
	imgURL, targetWidth, targetHeight, err := h.parseResizeIimageArgs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.app.GetResizedImage(h.ctx, r.Header, imgURL, targetWidth, targetHeight)
	if err != nil {
		h.logger.Error(
			"Resize error",
			zap.String("passed args", r.URL.Path),
			zap.String("requested URL", imgURL),
			zap.Int("width", targetWidth),
			zap.Int("heigth", targetHeight),
			zap.Error(err),
		)
		if errors.Is(err, imgproxy.ErrResponseNotOk) {
			http.Error(w, "image URL unavailable", http.StatusBadGateway)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	if _, writeErr := w.Write(result); writeErr != nil {
		h.logger.Error("Write response error", zap.Error(writeErr))
	}
}

func (h *Handler) parseResizeIimageArgs(requestPath string) (imageURL string, targetWidth int, targetHeight int, parseErr error) {
	// e.g /fill/300/200/url
	pathWithoutRoot := strings.Replace(requestPath, "/fill/", "", 1)
	pathArgs := strings.SplitN(pathWithoutRoot, "/", 3)

	if len(pathArgs) != 3 {
		h.logger.Error("Invalid path args", zap.Strings("parsed", pathArgs), zap.Error(ErrInvalidRequestParams))
		parseErr = ErrInvalidRequestParams
		return imageURL, targetWidth, targetHeight, parseErr
	}

	targetWidth, wErr := strconv.Atoi(pathArgs[0])
	if wErr != nil {
		h.logger.Error("Width convertation error", zap.String("parsed width", pathArgs[0]), zap.Error(wErr))
		parseErr = ErrInvalidTargetWidth
		return imageURL, targetWidth, targetHeight, parseErr
	}

	targetHeight, hErr := strconv.Atoi(pathArgs[1])
	if hErr != nil {
		h.logger.Error("Height convertation error", zap.String("parsed height", pathArgs[1]), zap.Error(hErr))
		parseErr = ErrInvalidTargetHeight
		return imageURL, targetWidth, targetHeight, parseErr
	}

	rawImageUrl := pathArgs[2]
	scheme := "https"
	if strings.HasPrefix(rawImageUrl, "http:/") {
		rawImageUrl = strings.Replace(rawImageUrl, "http:/", "", 1)
		scheme = "http"
	}

	imgURL, uErr := url.Parse(rawImageUrl)
	if uErr != nil {
		h.logger.Error("URL parse error", zap.String("parsed URL", pathArgs[2]), zap.Error(uErr))
		parseErr = ErrInvalidURL
		return imageURL, targetWidth, targetHeight, parseErr
	}
	if imgURL.Path == "" {
		h.logger.Error("empty URL", zap.String("parsed URL", pathArgs[2]))
		parseErr = ErrInvalidURL
		return imageURL, targetWidth, targetHeight, parseErr
	}

	if imgURL.Scheme == "" {
		imgURL.Scheme = scheme
	}

	return imgURL.String(), targetWidth, targetHeight, parseErr
}
