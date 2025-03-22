package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if !exists {
			// auth fail
			c.Abort()
		}
		c.Next()
	}
}
