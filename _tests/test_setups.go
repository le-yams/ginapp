package _tests

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
)

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
