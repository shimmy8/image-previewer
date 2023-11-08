package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Http  *HttpConfig
	Cache *CacheConfig
}

type HttpConfig struct {
	Port int `env:"HTTP_PORT" default:"8000"`
}

type CacheConfig struct {
	MaxSizeMB int `env:"CACHE_MAX_SIZE_MB" default:"256"`
}

func New() (*Config, error) {
	httpCnf := &HttpConfig{}
	if err := parseEnv(httpCnf); err != nil {
		return nil, err
	}

	cacheCnf := &CacheConfig{}
	if err := parseEnv(cacheCnf); err != nil {
		return nil, err
	}

	return &Config{Http: httpCnf, Cache: cacheCnf}, nil
}

func parseEnv(cnf interface{}) error {
	confPtr := reflect.ValueOf(cnf)
	ref := confPtr.Elem()

	a := ref.Kind()

	if a != reflect.Struct {
		return ErrInvalidConfig
	}

	for _, f := range reflect.VisibleFields(ref.Type()) {
		envName := f.Tag.Get("env")
		defaultVal := f.Tag.Get("default")
		envVal := getEnvValue(envName)

		value := envVal
		if envVal == "" {
			value = defaultVal
		}

		if f.Type.Kind() == reflect.Int {
			intVal, err := strconv.ParseInt(value, 10, 0)
			if err != nil {
				return err
			}
			ref.FieldByName(f.Name).SetInt(intVal)
		} else {
			ref.FieldByName(f.Name).SetString(value)
		}
	}
	return nil
}

func getEnvValue(envName string) string {
	for _, envStr := range os.Environ() {
		envSlice := strings.Split(envStr, "=")
		if len(envSlice) != 2 {
			continue
		}

		if envSlice[0] == envName {
			return envSlice[1]
		}
	}
	return ""
}
