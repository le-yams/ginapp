package ginapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartedAppHealthcheck(t *testing.T) {
	t.Parallel()
	// Arrange
	config := createTestConfig()
	config.Server.HealthcheckPath = "/healthz"

	app, err := New(&config, testSetups())
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

func (c TestConfig) GetServerConfig() *ServerConfig {
	return c.Server
}

func (c TestConfig) GetLogConfig() *LogConfig {
	return c.Log
}

func createTestConfig() TestConfig {
	return TestConfig{
		Server: &ServerConfig{
			Port: 0, // random port
			Mode: gin.TestMode,
		},
	}
}

type TestSetup struct {
	configureGinEngine func(*gin.Engine, *zap.SugaredLogger) error
	configureMetrics   func(*ginmetrics.Monitor) error
}

func testSetups() TestSetup {
	return TestSetup{
		configureGinEngine: func(engine *gin.Engine, logger *zap.SugaredLogger) error {
			return nil
		},
		configureMetrics: func(monitor *ginmetrics.Monitor) error {
			return nil
		},
	}
}

func (s TestSetup) ConfigureGinEngine(engine *gin.Engine, logger *zap.SugaredLogger) error {
	return s.configureGinEngine(engine, logger)
}

func (s TestSetup) ConfigureMetrics(monitor *ginmetrics.Monitor) error {
	return s.configureMetrics(monitor)
}
