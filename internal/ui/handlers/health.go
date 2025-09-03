package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmpl, err := template.ParseFiles("internal/ui/templates/health.html")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		err = tmpl.Execute(ctx.Writer, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}
}
