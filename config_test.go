package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	assertions "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
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

func TestWritableCustomConfig(t *testing.T) {
	t.Parallel()
	// Arrange
	config := createTestConfig()
	config.Custom = "initial"

	app, err := WithConfiguration(&config).
		ConfigureGinEngine(func(engine *gin.Engine, logger *zap.SugaredLogger) error {
			engine.GET("/custom-config", func(c *gin.Context) {
				cfg, _ := c.Get("config")
				customValue := cfg.(*testConfig).Custom
				cfg.(*testConfig).Custom = "changed"
				c.String(http.StatusOK, customValue)
			})
			return nil
		}).
		Build()
	if err != nil {
		t.Fatal(err)
	}

	// Act
	server := app.StartAsync()
	response, err := http.Get(fmt.Sprintf("http://%s/custom-config", server.Addr))
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(response.Body)
	_ = response.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = server.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assertions.Equal(t, "initial", string(body))
	assertions.Equal(t, "changed", config.Custom)
}
