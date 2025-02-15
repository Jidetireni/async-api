package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func NewLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		logger.Info("http request",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
		)

	}
}
