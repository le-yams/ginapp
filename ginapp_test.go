package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartedAppHealthcheck(t *testing.T) {
	t.Parallel()
	// Arrange
	config := &TestConfig{
		Server: &ServerConfig{
			Port:            0, // random port
			HealthcheckPath: "/healthz",
		},
		Log: &LogConfig{
			Level:  LogDebug,
			Format: LogJson,
		},
	}

	app, err := New(config, func(engine *gin.Engine, logger *zap.SugaredLogger) error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Act
	server := app.StartAsync()

	requestURL := fmt.Sprintf("http://%s%s", server.Addr, config.Server.HealthcheckPath)
	response, err := http.Get(requestURL)
	_ = response.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	err = server.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

type TestConfig struct {
	Server *ServerConfig
	Log    *LogConfig
}

func (c *TestConfig) GetServerConfig() *ServerConfig {
	return c.Server
}

func (c *TestConfig) GetLogConfig() *LogConfig {
	return c.Log
}
