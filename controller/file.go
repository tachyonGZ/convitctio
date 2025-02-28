package controller

import (
	"conviction/filesystem"
	"conviction/memocache"
	"conviction/model"
	"conviction/serializer"
	"conviction/util"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/coocood/freecache"
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

	// create placeholder
	head := filesystem.FileHead{
		MimeType:    param.MimeType,
		Name:        param.Name,
		Size:        param.Size,
		VirtualPath: param.Path,
	}

	if !fs.CreatePlaceHolder(&head) {
		c.String(500, "placeholder exists")
	}

	// store session in memocache
	uu, _ := uuid.NewRandom()
	uuM, _ := uu.MarshalText()
	key := string(uuM)
	uploadSession := serializer.UploadSession{
		Key:            key,
		UID:            fs.Owner.ID,
		VirtualPath:    head.VirtualPath,
		Name:           head.Name,
		Size:           head.Size,
		SavePath:       head.SavePath,
		LastModified:   head.LastModified,
		CallbackSecret: util.RandStringRunes(32),
	}
	memocache.SetJSON(append([]byte("callback_"), key...), uploadSession, ttl)

	// get credential
	credential := serializer.UploadCredential{
		SessionID: uploadSession.Key,
		Expires:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.JSON(200, credential)
}

func UploadBySession(c *gin.Context) {
	// binding
	/*
		var j struct {
			ID string `uri:"session_id" binding:"required"`
		}
	*/
	var param struct {
		ID string `json:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
	}

	// get upload session from cache
	value, err := memocache.GetJSON([]byte("callback_" + param.ID))
	if err != nil {
		c.String(500, err.Error())
	}
	uploadSession := value.(serializer.UploadSession)

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	/*
		file := filesystem.FileStream{
			File:        c.Request.Body,
			MimeType:    c.Request.Header.Get("Content-Type"),
			Name:        uploadSession.Name,
			Size:        uploadSession.Size,
			VirtualPath: uploadSession.VirtualPath,
		}
	*/

	head := filesystem.FileHead{
		MimeType:    c.Request.Header.Get("Content-Type"),
		Name:        uploadSession.Name,
		SavePath:    uploadSession.SavePath,
		Size:        uploadSession.Size,
		VirtualPath: uploadSession.VirtualPath,
	}

	fs.Upload(&head, c.Request.Body)
}

func CreateDownloadSession(c *gin.Context) {
	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// query object id from URL params
	objectID, _ := strconv.ParseUint(c.Query("object_id"), 10, 32)

	// get file by id
	file := model.GetFileByID(uint(objectID))

	// session
	// store file model on cache
	sessionID := util.RandStringRunes(16)
	cache, _ := c.Get("cache")
	b, _ := json.Marshal(file)
	cache.(*freecache.Cache).Set([]byte("download_"+sessionID), b, 60)

	// get download address
	downloadURL := fs.GetDownloadURL(uint(objectID), sessionID)

	c.JSON(0, downloadURL)
}

func DownloadBySession(c *gin.Context) {

	// binding
	var j struct {
		ID string `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(200, "")
	}

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// session
	// find file model on cache
	cache, _ := c.Get("cache")
	b, _ := cache.(*freecache.Cache).Get([]byte("download_" + j.ID))
	var target model.File
	json.Unmarshal(b, &target)

	// prepare for download
	rsc := fs.OpenFile(&target)
	defer fs.CloseFile(rsc)

	// send file
	c.Header("Content-Disposition", "attachment; filename=\""+url.PathEscape(target.Name)+"\"")
	http.ServeContent(c.Writer, c.Request, target.Name, target.UpdatedAt, rsc)
}

func Update(c *gin.Context) {

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// query object id from URL params
	objectID, _ := strconv.ParseUint(c.Query("object_id"), 10, 32)

	// get target
	target := model.GetFileByID(uint(objectID))

	f := filesystem.FileStream{}

	fs.UpdateFile(&target, f)
}

func Delete(c *gin.Context) {
	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))
}
