package middleware

import (
	"conviction/cache"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func TokenAuth(c *gin.Context, token string) int {
	user_uuid, err := cache.GetUserUUID(token)
	if err != nil {
		return 1
	}
	c.Set("user_id", user_uuid)
	return 0

}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		token := session.Get("token")
		code := TokenAuth(c, token.(string))

		if code != 0 {
			c.JSON(http.StatusUnauthorized, gin.H{})
			c.Abort()
		}

		c.Next()
	}
}
