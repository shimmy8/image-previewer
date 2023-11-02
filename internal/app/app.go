package app

import (
	"context"
	"image"

	"github.com/shimmy8/image-previewer/internal/config"
	"github.com/shimmy8/image-previewer/internal/service/cache"
	"github.com/shimmy8/image-previewer/internal/service/imgproxy"
	"github.com/shimmy8/image-previewer/internal/service/resizer"
)

type App struct {
	imgProxy *imgproxy.ImgProxy
	resizer  *resizer.Resizer
	cache    *cache.LruCache
}

func New(cacheConfig *config.CacheConfig) *App {
	return &App{
		imgProxy: imgproxy.New(),
		resizer:  resizer.New(),
		cache:    cache.New(cacheConfig),
	}
}

func (a *App) ResizeImage(ctx context.Context, url string, width int, heigth int) (image.Image, error) {
	// TODO
	return nil, nil
}
