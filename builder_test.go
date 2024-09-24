package ginapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Parallel()

	expectedError := "some error"
	t.Run("error configuring gin engine", func(t *testing.T) {
		t.Parallel()
		// Arrange
		config := createTestConfig()
		builder := WithConfiguration(&config).
			ConfigureGinEngine(func(engine *gin.Engine, logger *zap.SugaredLogger) error {
				return errors.New(expectedError)
			})

		// Act
		_, err := builder.Build()

		// Assert
		assert.EqualError(t, err, "error configuring gin engine: "+expectedError)
	})

	t.Run("error configuring metrics", func(t *testing.T) {
		t.Parallel()
		// Arrange
		config := createTestConfig()
		config.GetServerConfig().Metrics.Enabled = true

		builder := WithConfiguration(&config).
			ConfigureMetrics(func(metrics *ginmetrics.Monitor) error {
				return errors.New(expectedError)
			})

		// Act
		_, err := builder.Build()

		// Assert
		assert.EqualError(t, err, "error configuring metrics: "+expectedError)
	})
}
