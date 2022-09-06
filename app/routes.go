package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) addV1Routes(r *gin.Engine) {
	apiV1 := r.Group("/api/v1")

	apiV1.GET("/ping", jsonHandler(pingHandler))

	a.addRunRoutes(apiV1)
	a.addJobRoutes(apiV1)
}

func pingHandler(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, gin.H{"message": "pong"}, nil
}
