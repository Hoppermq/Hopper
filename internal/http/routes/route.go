// Package routes represents the http routes package.
package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterBaseRoutes register http routes
func RegisterBaseRoutes(e *gin.Engine) {
	e.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "hello world")
	})
}
