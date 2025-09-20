package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) Middleware {
	return Middleware{logger}
}

// Logger logs each incoming HTTP request and response status
func (m Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		m.logger.Info("request",
			slog.Int("status", status),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("query", c.Request.URL.RawQuery),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("latency", latency.String()),
		)
	}
}

// Recovery handles panics and logs them
func (m Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				m.logger.Error("panic recovered",
					slog.Any("error", rec),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
