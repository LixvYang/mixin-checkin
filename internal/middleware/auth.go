package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Auth(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("auth-id") != "" {
			c.Abort()
		}

		c.Next()
	}
}
