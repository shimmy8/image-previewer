package config

type Config struct {
	Http  HttpConfig
	Cache CacheConfig
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT" default:"8000"`
}

type CacheConfig struct {
	MaxSizeMB int `env:"CACHE_MAX_SIZE_MB" default:"256"`
}

func New() *Config {
	// TODO
	return &Config{}
}
