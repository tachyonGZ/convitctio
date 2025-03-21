package controller

import (
	"conviction/filesystem"
	"conviction/model"

	"github.com/gin-gonic/gin"
)

func CreateSharedFile(c *gin.Context) {
	// check binding
	var param struct {
		SourceID string `json:"source_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(200, "bad binding"+err.Error())
		return
	}

	// create file system
	u, exists := c.Get("user")
	if !exists {
		c.String(500, "current user not exists")
		return
	}
	fs := filesystem.NewFileSystem(u.(*model.User))

	// create shared file
	sharedFileID, e := fs.CreateSharedFile(param.SourceID)
	if e != nil {
		c.String(500, "create shared file fail")
		return
	}

	// response
	c.JSON(
		200,
		struct {
			SharedFileID string `json:"shared_file_id"`
		}{
			SharedFileID: sharedFileID,
		})
}

func DeleteSharedFile(c *gin.Context) {
	// check binding
	var param struct {
		SharedFileID string `json:"shared_file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(200, "bad binding"+err.Error())
		return
	}

	// create file system
	u, exists := c.Get("user")
	if !exists {
		c.String(500, "current user not exists")
		return
	}
	fs := filesystem.NewFileSystem(u.(*model.User))

	fs.DeleteSharedFile(param.SharedFileID)

	c.String(200, "")
}

func CreateSharedDownloadSession(c *gin.Context) {

}
