package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	HTTP  *HTTPConfig
	Cache *CacheConfig
	Proxy *ProxyConfig
}

type HTTPConfig struct {
	Port int `env:"HTTP_PORT" default:"8000"`
}

type CacheConfig struct {
	MaxElemCnt int    `env:"CACHE_MAX_ELEM_CNT" default:"50"`
	Dir        string `env:"CACHE_DIR" default:"./filecache"`
}

type ProxyConfig struct {
	Timeout int `env:"PROXY_TIMEOUT" default:"2"`
}

func New() (*Config, error) {
	httpCnf := &HTTPConfig{}
	if err := parseEnv(httpCnf); err != nil {
		return nil, err
	}

	cacheCnf := &CacheConfig{}
	if err := parseEnv(cacheCnf); err != nil {
		return nil, err
	}

	proxyCnf := &ProxyConfig{}
	if err := parseEnv(proxyCnf); err != nil {
		return nil, err
	}

	return &Config{HTTP: httpCnf, Cache: cacheCnf, Proxy: proxyCnf}, nil
}

func parseEnv(cnf interface{}) error {
	confPtr := reflect.ValueOf(cnf)
	ref := confPtr.Elem()

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
