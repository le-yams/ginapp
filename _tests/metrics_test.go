package _tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/le-yams/ginapp"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"

	assertions "github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	t.Parallel()

	t.Run("disabled", func(t *testing.T) {
		t.Parallel()
		config := createTestConfig()

		app, err := ginapp.New(&config, testSetups())
		if err != nil {
			t.Fatal(err)
		}

		// Act
		server := app.StartAsync()

		requestURL := fmt.Sprintf("http://%s/metrics", server.Addr)
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
		assertions.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("enabled with custom counter", func(t *testing.T) {
		t.Parallel()
		metricsPath := "/my-metrics"
		config := createTestConfig()
		config.Server.Metrics.Enabled = true
		config.Server.Metrics.Path = metricsPath
		metricName := "test_count"
		setups := testSetups()
		setups.configureMetrics = func(monitor *ginmetrics.Monitor) error {
			testMetric := &ginmetrics.Metric{
				Type:        ginmetrics.Counter,
				Name:        metricName,
				Description: "Test counter",
				Labels:      []string{},
			}
			return ginmetrics.GetMonitor().AddMetric(testMetric)
		}
		setups.configureGinEngine = func(engine *gin.Engine, logger *zap.SugaredLogger) error {
			engine.GET("/test", func(c *gin.Context) {
				metric := ginmetrics.GetMonitor().GetMetric(metricName)
				err := metric.Inc([]string{})
				if err != nil {
					t.Fatal(err)
				}
				c.Status(http.StatusOK)
			})
			return nil
		}
		app, err := ginapp.New(&config, setups)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		server := app.StartAsync()
		metricsURL := fmt.Sprintf("http://%s%s", server.Addr, metricsPath)
		response, err := http.Get(metricsURL)
		if err != nil {
			t.Fatal(err)
		}
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		metrics1 := string(bodyBytes)
		_ = response.Body.Close()

		response, err = http.Get(fmt.Sprintf("http://%s/test", server.Addr))
		if err != nil {
			t.Fatal(err)
		}
		_ = response.Body.Close()

		response, err = http.Get(metricsURL)
		if err != nil {
			t.Fatal(err)
		}
		bodyBytes, err = io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		metrics2 := string(bodyBytes)
		_ = response.Body.Close()

		if err != nil {
			t.Fatal(err)
		}
		err = server.Close()
		if err != nil {
			t.Fatal(err)
		}

		// Assert
		assert := assertions.New(t)
		assert.Equal(http.StatusOK, response.StatusCode)
		assert.NotContains(metrics1, metricName)
		assert.Contains(metrics2, metricName+" 1")
	})
}
