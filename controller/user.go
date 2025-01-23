package controller

import (
	"conviction/db"
	"conviction/model"
	"conviction/serializer"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {

	// binding
	var param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, serializer.Response{})
	}

	u, err := model.GetUserByUsername(param.Username)

	// auth username
	if nil != err {
		c.String(200, "Wrong password or email address")
		fmt.Println("Wrong password or email address")
		return
	}

	// auth password
	if authOK := u.CheckPassword(param.Password); !authOK {
		c.String(200, "Wrong password or email address")
		fmt.Println("Wrong password or email address")
		return
	}

	// session
	session := sessions.Default(c)
	session.Set("user_id", u.ID)

	c.JSON(200, serializer.Response{})
}

func UserRegister(c *gin.Context) {

	// binding
	var param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, serializer.Response{})
	}

	// create new user on db
	u := model.NewUser()
	u.Username = param.Username
	u.Password = param.Password
	db.GetDB().Create(&u)

	c.JSON(200, serializer.Response{})
}
