package controller

import (
	"conviction/model"
	"conviction/serializer"

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

	user, err := model.FindUserByUsername(param.Username)

	// auth username
	if nil != err {
		c.String(200, "Wrong password or email address")
		return
	}

	// auth password
	if authOK := user.CheckPassword(param.Password); !authOK {
		c.String(200, "Wrong password or email address")
		return
	}

	// session
	session := sessions.Default(c)
	session.Set("user_id", user.UUID)
	session.Save()

	// get user root dir
	rootID, _ := model.GetUserRootID(user.UUID)

	// response
	c.JSON(
		200,
		struct {
			RootID string `json:"root_id"`
		}{
			RootID: rootID,
		})
}

func UserRegister(c *gin.Context) {

	// datga binding
	var param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
		return
	}

	// create user
	user := model.User{
		Username: param.Username,
		Password: param.Password,
	}
	res := user.Create()
	if res != nil {
		c.String(500, res.Error())
		return
	}

	c.JSON(200, serializer.Response{})
}
