package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hoppermq/hopper/internal/ui/handlers"
)

// RegisterBaseRoutes register to the engine the routes under "/".
func RegisterBaseRoutes(e *gin.Engine) {
	e.GET("/", handlers.DashboardHandler())
}
