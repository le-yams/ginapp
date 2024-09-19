package ginapp

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	t.Parallel()
	// Arrange
	config := createTestConfig()
	customHealthCheckPath := "/healthz"
	config.Server.HealthcheckPath = customHealthCheckPath

	app, err := WithConfiguration(&config).Build()
	if err != nil {
		t.Fatal(err)
	}

	// Act
	server := app.StartAsync()

	requestURL := fmt.Sprintf("http://%s%s", server.Addr, customHealthCheckPath)
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
