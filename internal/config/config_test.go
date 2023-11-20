package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("test parse env", func(t *testing.T) {
		testPort := 777
		os.Setenv("HTTP_PORT", strconv.Itoa(testPort))

		testCacheSize := 50
		os.Setenv("CACHE_MAX_ELEM_CNT", strconv.Itoa(testCacheSize))

		testCacheDir := "/tmp/cache"
		os.Setenv("CACHE_DIR", testCacheDir)

		testProxyTimeout := 12
		os.Setenv("PROXY_TIMEOUT", strconv.Itoa(testProxyTimeout))

		cnf, err := New()
		require.NoError(t, err)

		require.Equal(t, testPort, cnf.HTTP.Port)
		require.Equal(t, testCacheSize, cnf.Cache.MaxElemCnt)
		require.Equal(t, testCacheDir, cnf.Cache.Dir)
		require.Equal(t, testProxyTimeout, cnf.Proxy.Timeout)
	})
}
