package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")

		if nil == user {
			fmt.Println("not login")
			c.Abort()
		}

		c.Next()
	}
}
