package ginapp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupHealthcheck(ginEngine *gin.Engine, serverConfig *ServerConfig) {
	ginEngine.GET(serverConfig.GetHealthCheckPath(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
