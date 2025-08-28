package routes

import "github.com/gin-gonic/gin"

func RegisterBaseRoutes(e *gin.Engine) {
	e.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})
}
