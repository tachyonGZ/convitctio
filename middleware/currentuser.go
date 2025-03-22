package middleware

import (
	"conviction/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		user_id := session.Get("user_id")
		if nil == user_id {
			// not ser user
			c.Next()
		}

		c.Set("user_id", user_id)

		user, err := model.FindUser(user_id.(string))
		if err != nil {
			c.Next()
		}

		c.Set("user", &user)
		c.Next()
	}
}
