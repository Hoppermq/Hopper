package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tmpl, err := template.ParseFiles("internal/ui/templates/login.html")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = tmpl.Execute(ctx.Writer, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}

