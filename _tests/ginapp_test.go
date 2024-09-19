package _tests

import (
	"fmt"
	"github.com/le-yams/ginapp"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartedAppHealthcheck(t *testing.T) {
	t.Parallel()
	// Arrange
	config := createTestConfig()
	config.Server.HealthcheckPath = "/healthz"

	app, err := ginapp.New(&config, testSetups())
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
