package controller

import (
	"conviction/filesystem"
	"conviction/memocache"
	"conviction/serializer"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

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
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	fs.DeleteSharedFile(param.SharedFileID)

	c.String(200, "")
}

func CreateSharedFileDownloadSession(c *gin.Context) {

	ttl := 6000

	// data binding
	var param struct {
		SharedFileID string `json:"shared_file_id" binding:"required"`
	}
	if e := c.ShouldBindJSON(&param); e != nil {
		c.JSON(500, e.Error())
		return
	}

	// store session in cache
	key, e := uuid.NewRandom()
	if e != nil {
		c.JSON(500, e.Error())
	}
	session := serializer.DownloadSession{
		Key: key.String(),

		DestType: serializer.SharedFile,
		DestID:   param.SharedFileID,
	}
	memocache.SetDownloadSession(key.String(), &session, ttl)

	// get credential
	credential := serializer.DownloadCredential{
		SessionID: session.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}
