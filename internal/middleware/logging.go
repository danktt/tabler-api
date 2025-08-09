package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging(logger *slog.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        status := c.Writer.Status()
        logger.Info(
            "handled request",
            slog.Int("statusCode", status),
            slog.String("method", c.Request.Method),
            slog.String("path", c.Request.URL.Path),
            slog.Any("duration", time.Since(start)),
        )
    }
}
