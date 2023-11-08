package cache

import "github.com/shimmy8/image-previewer/internal/config"

type LruCache struct {
}

func New(conf *config.CacheConfig) *LruCache {
	//TODO
	return &LruCache{}
}

func (c *LruCache) Get(key string) ([]byte, error) {
	return nil, nil
}

func (c *LruCache) Set(key string, value []byte) error {
	return nil
}
