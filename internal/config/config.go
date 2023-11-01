package config

type Config struct {
	Http  HttpConfig
	Cache CacheConfig
}

type HttpConfig struct {
	Port int `env:"IMG_PREVIEWER_HTTP_PORT" default:"8000"`
}

type CacheConfig struct {
	MaxSizeBytes int `env:"IMG_PREVIEWER_CACHE_MAX_SIZE_BYTES" default:"1024"`
	MaxElements  int `env:"IMG_PREVIEWER_CACHE_MAX_ELEMENTS" default:"100"`
}

func New() *Config {
	// TODO
	return &Config{}
}
