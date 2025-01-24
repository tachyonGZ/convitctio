package controller

import (
	"conviction/filesystem"
	"conviction/model"
	"conviction/serializer"
	"conviction/util"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUploadSession(c *gin.Context) {

	var ttl int = 6000

	// check binding
	var param struct {
		Path string `json:"path" binding:"required"`
		Size uint64 `json:"size" binding:"min=0"`
		Name string `json:"name" binding:"required"`
		//PolicyID     string `json:"policy_id" binding:"required"`
		LastModified int64  `json:"last_modified"`
		MimeType     string `json:"mime_type"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(200, "bad binding")
		fmt.Println("bad binding")
	}

	// create file system
	u, exists := c.Get("user")
	if !exists {
		fmt.Println("user not exists")
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
	fs.CreatePlaceHolder(&head)

	// store session in cache
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
	b, _ := json.Marshal(uploadSession)
	cache, _ := c.Get("cache")
	cache.(*freecache.Cache).Set(append([]byte("callback_"), key...), b, ttl)

	// get credential
	credential := fs.Adapter.Token(uploadSession)

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
		c.JSON(200, "")
	}

	// get upload session from cache
	cache, _ := c.Get("cache")
	uploadSessionRaw, _ := cache.(*freecache.Cache).Get([]byte("callback_" + param.ID))
	var uploadSession serializer.UploadSession
	json.Unmarshal(uploadSessionRaw, &uploadSession)

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
