package controller

import (
	"conviction/filesystem"
	"conviction/memocache"
	"conviction/model"
	"conviction/serializer"
	"conviction/util"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUploadSession(c *gin.Context) {

	var ttl int = 6000
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
	u, exists := c.Get("user")
	if !exists {
		c.String(500, "current user not exists")
		return
	}
	fs := filesystem.NewFileSystem(u.(*model.User))

	// open or create directory
	dir_id, exist, _ := fs.OpenDirectory(param.Path)
	if !exist {
		dir_id = fs.CreateDirectoryByPath(param.Path)
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
	memocache.SetUploadSession(key, &uploadSession, ttl)

	// get credential
	credential := serializer.UploadCredential{
		SessionID: uploadSession.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}

func UploadBySession(c *gin.Context) {
	// data binding
	var param struct {
		SessionID string `uri:"session_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(500, err.Error())
		return
	}

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// get upload session from cache
	pSession, err := memocache.GetUploadSession(param.SessionID)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	head := filesystem.FileHead{
		MimeType:    pSession.MimeType,
		Name:        pSession.Name,
		SavePath:    pSession.SavePath,
		Size:        pSession.Size,
		VirtualPath: pSession.VirtualPath,
	}

	fs.Upload(&head, c.Request.Body, pSession.PlaceholderID)

	// delete upload session in cache
	memocache.DeleteUploadSession(pSession.Key)

	c.String(200, "")
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

	// get current user
	u, _ := c.Get("user")
	user := u.(*model.User)

	// create file system
	fs := filesystem.NewFileSystem(user)

	fid, _ := strconv.ParseUint(param.FileID, 10, 32)
	head := fs.GetFileHead(uint(fid))

	// store session in cache
	uu, _ := uuid.NewRandom()
	uuM, _ := uu.MarshalText()
	key := string(uuM)
	session := serializer.DownloadSession{
		Key: key,

		FileID:  param.FileID,
		Name:    head.Name,
		OwnerID: fs.Owner.UUID,
	}
	memocache.SetDownloadSession(key, &session, ttl)

	// get credential
	credential := serializer.DownloadCredential{
		SessionID: session.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}

func Update(c *gin.Context) {

	u, _ := c.Get("user")
	user := u.(*model.User)

	// create file system
	fs := filesystem.NewFileSystem(user)

	// query object id from URL params
	objectID, _ := strconv.ParseUint(c.Query("object_id"), 10, 32)

	// get target
	target, _ := model.GetFileByID(uint(objectID), user.ID)

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

	// get user
	u, _ := c.Get("user")
	user := u.(*model.User)

	// create file system
	fs := filesystem.NewFileSystem(user)

	// delete file
	if err := fs.DeleteFile(param.FileID); err != nil {
		c.String(500, err.Error())
		return
	}

	// response
	c.String(200, "")
}

func GetFileInfo(c *gin.Context) {
	// data binding
	var param struct {
		FileID string `json:"file_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(500, err.Error())
		return
	}
}
