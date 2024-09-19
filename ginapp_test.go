package ginapp

import (
	"github.com/gin-gonic/gin"
)

type testConfig struct {
	Server *ServerConfig
	Log    *LogConfig
	Custom string
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
