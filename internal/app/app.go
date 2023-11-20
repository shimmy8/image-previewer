package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shimmy8/image-previewer/internal/service/cache"
	"github.com/shimmy8/image-previewer/internal/service/imgproxy"
	"github.com/shimmy8/image-previewer/internal/service/resizer"
	"go.uber.org/zap"
)

type App struct {
	imgProxy *imgproxy.ImgProxy
	resizer  *resizer.Resizer
	cache    *cache.LruCache

	logger *zap.Logger
}

func New(cacheMaxSize int, cacheDir string, proxyTimeout int, logger *zap.Logger) *App {
	return &App{
		imgProxy: imgproxy.New(proxyTimeout),
		resizer:  resizer.New(),
		cache:    cache.New(cacheMaxSize, cacheDir),
		logger:   logger,
	}
}

func (a *App) GetResizedImage(
	ctx context.Context,
	headers map[string][]string,
	url string,
	width int,
	heigth int,
) ([]byte, error) {
	imgCacheKey := strings.Join([]string{strconv.Itoa(width), strconv.Itoa(heigth), url}, ":")

	cachedResult, err := a.cache.Get(imgCacheKey)
	if err != nil {
		a.logger.Named("cache").Error("Get from cache", zap.String("key", imgCacheKey), zap.Error(err))
	}

	if cachedResult != nil {
		return cachedResult, nil
	}

	originalImg, err := a.imgProxy.GetImage(ctx, url, headers)
	if err != nil {
		a.logger.Named("imgProxy").Error("Get from url", zap.String("url", url), zap.Error(err))
		if errors.Is(err, imgproxy.ErrResponseNotOk) {
			return nil, fmt.Errorf("%w: %w", ErrProxyResponseNotOk, err)
		}
		return nil, fmt.Errorf("%w: %w", ErrProxyGetImage, err)
	}

	resizedImg, err := a.resizer.ResizeImage(originalImg, width, heigth)
	if err != nil {
		a.logger.Named("resizer").Error("Resize", zap.Error(err))
		switch {
		case errors.Is(err, resizer.ErrFormatNotSupported):
			return nil, fmt.Errorf("%w: %w", ErrImageFormat, err)
		case errors.Is(err, resizer.ErrNotAnImage):
			return nil, fmt.Errorf("%w: %w", ErrFileNotAnImage, err)
		default:
			return nil, fmt.Errorf("%w: %w", ErrResizeImage, err)
		}
	}

	cacheSetErr := a.cache.Set(imgCacheKey, resizedImg)
	if cacheSetErr != nil {
		a.logger.Named("cache").Error("Set img to cache", zap.String("key", imgCacheKey), zap.Error(cacheSetErr))
	}

	return resizedImg, nil
}
