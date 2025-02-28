package middleware

import (
	"conviction/model"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		uid := session.Get("user_id")
		if nil == uid {
			fmt.Println("session not set user_id")
			c.Next()
		}

		user, err := model.GetUserByID(uid)
		if err != nil {
			c.Next()
		}

		c.Set("user", &user)
		c.Next()
	}
}
