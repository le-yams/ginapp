package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
)

type testConfig struct {
	Server *ServerConfig
	Log    *LogConfig
}

func (c testConfig) GetServerConfig() *ServerConfig {
	return c.Server
}

func (c testConfig) GetLogConfig() *LogConfig {
	return c.Log
}

func createTestConfig() testConfig {
	return testConfig{
		Server: &ServerConfig{
			Port:    0, // random port
			Mode:    gin.TestMode,
			Metrics: &MetricsConfig{},
		},
	}
}

type setups struct {
	configureGinEngine func(*gin.Engine, *zap.SugaredLogger) error
	configureMetrics   func(*ginmetrics.Monitor) error
}

func testSetups() setups {
	return setups{
		configureGinEngine: func(engine *gin.Engine, logger *zap.SugaredLogger) error {
			return nil
		},
		configureMetrics: func(monitor *ginmetrics.Monitor) error {
			return nil
		},
	}
}

func (s setups) ConfigureGinEngine(engine *gin.Engine, logger *zap.SugaredLogger) error {
	return s.configureGinEngine(engine, logger)
}

func (s setups) ConfigureMetrics(monitor *ginmetrics.Monitor) error {
	return s.configureMetrics(monitor)
}
