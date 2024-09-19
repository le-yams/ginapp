package _tests

import (
	"github.com/gin-gonic/gin"
	"github.com/le-yams/ginapp"
)

type testConfig struct {
	Server *ginapp.ServerConfig
	Log    *ginapp.LogConfig
}

func (c testConfig) GetServerConfig() *ginapp.ServerConfig {
	return c.Server
}

func (c testConfig) GetLogConfig() *ginapp.LogConfig {
	return c.Log
}

func createTestConfig() testConfig {
	return testConfig{
		Server: &ginapp.ServerConfig{
			Port:    0, // random port
			Mode:    gin.TestMode,
			Metrics: &ginapp.MetricsConfig{},
		},
	}
}
