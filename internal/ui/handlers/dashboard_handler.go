package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tmpl, err := template.ParseFiles("internal/ui/templates/dashboard.html")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error parsing template"})
			return
		}

		err = tmpl.Execute(ctx.Writer, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error executing template"})
			return
		}
	}
}
