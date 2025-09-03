package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")

		ctx.Status(http.StatusOK)
		i := 0
		t  := []string{"ping", "pong"}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

				if _, err := fmt.Fprintf(ctx.Writer, "event: ping\ndata: %v\n\n", t[i]); err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if f, ok := ctx.Writer.(http.Flusher); ok {
					f.Flush()
				}

				i = (i + 1) % len(t)
				time.Sleep(time.Second * 2)
			}
	}
}
