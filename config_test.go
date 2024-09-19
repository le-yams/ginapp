package ginapp

import (
	"github.com/gin-gonic/gin"
	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultServerConfig(t *testing.T) {
	t.Parallel()
	t.Run("default server configuration", func(t *testing.T) {
		t.Parallel()
		serverConfig := defaultServerConfig()

		assert := assertions.New(t)
		assert.Equal(0, serverConfig.Port)
		assert.Equal(gin.ReleaseMode, serverConfig.Mode)
		assert.Equal("/health", serverConfig.HealthcheckPath)

		metrics := serverConfig.Metrics
		require.NotNil(t, metrics)
		assert.False(metrics.Enabled)
		assert.Equal("/metrics", metrics.Path)
	})

	t.Run("default healthcheck path", func(t *testing.T) {
		t.Parallel()
		assert := assertions.New(t)
		assert.Equal("/health", (&ServerConfig{}).GetHealthCheckPath())
	})

	t.Run("default metrics configuration", func(t *testing.T) {
		t.Parallel()
		metrics := (&ServerConfig{}).GetMetrics()
		assertions.NotNil(t, metrics)
		assertions.Equal(t, "/metrics", metrics.Path)
	})

	t.Run("default metrics path", func(t *testing.T) {
		t.Parallel()
		assertions.Equal(t, "/metrics", (&MetricsConfig{}).GetPath())
	})
}
