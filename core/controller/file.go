package controller

import (
	"conviction/cache"
	"conviction/filesystem"
	"conviction/model"
	"conviction/serializer"
	"conviction/util"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUploadSession(c *gin.Context) {

	var ttl int64 = 6000
	if os.Getenv("DEBUG") != "" {
		ttl = 0
	}
	// check binding
	var param struct {
		LastModified int64  `json:"last_modified"`
		MimeType     string `json:"mime_type" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Path         string `json:"path" binding:"required"`
		Size         uint64 `json:"size" binding:"min=0"`
		//PolicyID     string `json:"policy_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(200, "bad binding"+err.Error())
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// open or create directory
	dir_id, exist, e := fs.OpenDirectory(param.Path)
	if e != nil {
		c.String(500, e.Error())
		return
	}
	if !exist {
		dir_id, e = fs.CreateDirectoryByPath(param.Path)
		if e != nil {
			c.String(500, e.Error())
			return
		}
	}

	head := filesystem.FileHead{
		MimeType:    param.MimeType,
		Name:        param.Name,
		Size:        param.Size,
		VirtualPath: param.Path,
	}

	// create placeholder
	holder_id, err := fs.CreatePlaceHolder(&head, dir_id)
	if err != nil {
		c.String(500, "create placeholder fail"+err.Error())
		return
	}

	// store session in memocache
	uu, _ := uuid.NewRandom()
	uuM, _ := uu.MarshalText()
	key := string(uuM)
	uploadSession := serializer.UploadSession{
		Key: key,

		PlaceholderID:  holder_id,
		OwnerID:        fs.Owner.UUID,
		VirtualPath:    head.VirtualPath,
		MimeType:       head.MimeType,
		Name:           head.Name,
		Size:           head.Size,
		SavePath:       head.SavePath,
		LastModified:   head.LastModified,
		CallbackSecret: util.RandStringRunes(32),
	}
	cache.SetUploadSession(key, &uploadSession, ttl)

	// get credential
	credential := serializer.UploadCredential{
		SessionID: uploadSession.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}

func CreateDownloadSession(c *gin.Context) {

	ttl := 60

	// data binding
	var param struct {
		FileID string `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
		return
	}

	// store session in cache
	uu, _ := uuid.NewRandom()
	uuM, _ := uu.MarshalText()
	key := string(uuM)
	session := serializer.DownloadSession{
		Key: key,

		DestType: serializer.PersonalFile,
		DestID:   param.FileID,
	}
	cache.SetDownloadSession(key, &session, int64(ttl))

	// get credential
	credential := serializer.DownloadCredential{
		SessionID: session.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}

func Update(c *gin.Context) {

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// query object id from URL params
	object_id := c.Query("object_id")

	// get target
	target, _ := model.FindUserFile(user_id.(string), object_id)

	f := filesystem.FileStream{}

	fs.UpdateFile(target, f)
}

func DeleteFile(c *gin.Context) {
	// data binding
	var param struct {
		FileID string `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(500, err.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// delete file
	if err := fs.DeleteFile(param.FileID); err != nil {
		c.String(500, err.Error())
		return
	}

	// response
	c.String(200, "")
}

func GetFileStatus(c *gin.Context) {
	// data binding
	var param struct {
		FileID string `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(500, err.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	p_head, e := fs.GetFileHead(param.FileID)
	if e != nil {
		c.String(500, e.Error())
		return
	}

	// response
	c.JSON(
		200,
		struct {
			Name     string `json:"name"`
			MimeType string `json:"mime_type"`
			Size     uint64 `json:"size"`
		}{
			Name:     p_head.Name,
			MimeType: p_head.MimeType,
			Size:     p_head.Size,
		})
}

func MoveFile(c *gin.Context) {
}

func RenameFile(c *gin.Context) {
	// data binding
	var param struct {
		FileID string `json:"file_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
	}
	if e := c.ShouldBindJSON(&param); e != nil {
		c.String(500, e.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	e := fs.RenameFile(param.FileID, param.Name)
	if e != nil {
		c.String(500, e.Error())
		return
	}

	c.String(200, "")
}
