package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

func GinXid(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := xid.New().String()
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("xid", correlationID)
		})
		c.Header("X-Correlation-ID", correlationID)
		c.Set("xid", logger)

		c.Next()
	}
}
