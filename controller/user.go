package controller

import (
	"conviction/model"
	"conviction/serializer"
	"fmt"
	"strconv"

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
	session.Save()
	fmt.Println("UserLogin: set session")

	// get user root dir
	rootID, _ := model.GetUserRootID(u.ID)
	rootIDRaw := strconv.FormatUint(uint64(rootID), 10)

	// response
	c.JSON(
		200,
		struct {
			RootID string `json:"root_id"`
		}{
			RootID: rootIDRaw,
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
