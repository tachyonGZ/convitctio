package controller

import (
	"conviction/db"
	"conviction/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	user := model.GetUserByID()

	session := sessions.Default(c)

	session.Set("user_id", user.ID)

}

func UserRegister(c *gin.Context) {

	// binding
	var param struct {
		Username string `form:"userName" json:"userName" binding:"required,email"`
		Password string `form:"Password" json:"Password" binding:"required,min=4,max=64"`
	}

	// create new user
	u := model.NewUser()
	u.Username = param.Username
	u.Password = param.Password

	db.GetDB().Create(&u)
}
