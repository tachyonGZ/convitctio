package middlewware

import "github.com/gin-gonic/gin"

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")

		if nil == user {
			c.Abort()
		}

		c.Next()
	}
}
