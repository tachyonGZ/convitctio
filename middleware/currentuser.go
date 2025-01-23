package middlewware

import (
	"conviction/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		uid := session.Get("user_id")
		if nil == uid {
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
