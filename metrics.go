package ginapp

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func setupMetrics(ginEngine *gin.Engine, configuration *MetricsConfig, configure func(*ginmetrics.Monitor) error) error {
	if !configuration.Enabled {
		return nil
	}

	monitor := ginmetrics.GetMonitor()
	monitor.SetMetricPath(configuration.GetPath())
	if configure != nil {
		err := configure(monitor)
		if err != nil {
			return err
		}
	}
	monitor.Use(ginEngine)
	return nil
}
