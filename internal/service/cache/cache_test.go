package cache

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shimmy8/image-previewer/internal/config"
	"github.com/stretchr/testify/require"
)

func TestCacheGetSet(t *testing.T) {
	cacheFolder := "./testcache"
	diskCache := New(&config.CacheConfig{MaxSize: 4, Dir: cacheFolder})
	cacheBytes := []byte("some")

	defer func() { os.RemoveAll(cacheFolder) }()

	t.Run("test cache folder created", func(t *testing.T) {
		_, err := os.Stat(cacheFolder)
		require.NoError(t, err)
	})

	t.Run("test set - get", func(t *testing.T) {
		err := diskCache.Set("test", cacheBytes)

		require.NoError(t, err)

		cached, err := diskCache.Get("test")
		require.NoError(t, err)

		require.Equal(t, cacheBytes, cached)
	})

	t.Run("test cache max size", func(t *testing.T) {
		cacheFolder2 := "./testcache2"
		newCache := New(&config.CacheConfig{MaxSize: 2, Dir: cacheFolder2})

		defer func() { os.RemoveAll(cacheFolder2) }()

		key1 := "test_1"
		val1 := []byte("one")

		key2 := "test_2"
		val2 := []byte("more")

		key3 := "test_3"
		val3 := []byte("thing")

		// add val1 to cache
		err1 := newCache.Set(key1, val1)
		require.NoError(t, err1)
		// add val2 to cache
		err2 := newCache.Set(key2, val2)
		require.NoError(t, err2)
		// test val1 in cache
		cached1, errG1 := newCache.Get(key1)
		require.NoError(t, errG1)
		require.Equal(t, val1, cached1)
		// test val2 in cache
		cached2, errG2 := newCache.Get(key2)
		require.NoError(t, errG2)
		require.Equal(t, val2, cached2)
		// add val3 to cache
		err3 := newCache.Set(key3, val3)
		require.NoError(t, err3)
		// test val3 in cache
		cached3, errG3 := newCache.Get(key3)
		require.NoError(t, errG3)
		require.Equal(t, val3, cached3)

		// test val1 not in cache
		uncached1, errG1 := newCache.Get(key1)
		require.NoError(t, errG1)
		require.Nil(t, uncached1)
	})
}

func TestCahceLoad(t *testing.T) {
	cacheFolder := "./testcache3"
	defer func() { os.RemoveAll(cacheFolder) }()

	testKey := "test"
	testValue := []byte("value")

	filename := filepath.Join(cacheFolder, base64.StdEncoding.EncodeToString([]byte(testKey)))
	mkErr := os.MkdirAll(cacheFolder, os.ModePerm)
	require.NoError(t, mkErr)

	writeErr := os.WriteFile(filename, testValue, os.ModePerm)
	require.NoError(t, writeErr)

	t.Run("test cache load from folder", func(t *testing.T) {
		diskCache := New(&config.CacheConfig{MaxSize: 2, Dir: cacheFolder})

		time.Sleep(time.Millisecond * 1)

		val, err := diskCache.Get(testKey)
		require.NoError(t, err)
		require.Equal(t, testValue, val)
	})
}
