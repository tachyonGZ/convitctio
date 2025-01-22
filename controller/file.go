package controller

import (
	"conviction/filesystem"
	"conviction/model"
	"conviction/serializer"
	"conviction/util"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
)

func CreateUploadSession(c *gin.Context) {

	// check binding
	var j struct {
		Path string `json:"path" binding:"required"`
		Size uint64 `json:"size" binding:"min=0"`
		Name string `json:"name" binding:"required"`
		//PolicyID     string `json:"policy_id" binding:"required"`
		LastModified int64  `json:"last_modified"`
		MimeType     string `json:"mime_type"`
	}
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(200, "")
	}

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	file := filesystem.FileStream{
		File:        io.NopCloser(strings.NewReader("")),
		MimeType:    j.MimeType,
		Name:        j.Name,
		Size:        j.Size,
		VirtualPath: j.Path,
	}

	var callbackKey string

	//
	uploadSession := serializer.UploadSession{
		Key:            callbackKey,
		UID:            fs.Owner.ID,
		VirtualPath:    file.VirtualPath,
		Name:           file.Name,
		Size:           file.Size,
		SavePath:       file.SavePath,
		LastModified:   file.LastModified,
		CallbackSecret: util.RandStringRunes(32),
	}

	// get credential
	credential := fs.Adapter.Token(uploadSession)
	// create placeholder

	// session
	b, _ := json.Marshal(uploadSession)
	cache, _ := c.Get("cache")
	cache.(*freecache.Cache).Set([]byte("callback_"+callbackKey), b, 60)

	c.JSON(0, credential)
}

func UploadBySession(c *gin.Context) {
	// binding
	var j struct {
		ID string `uri:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(200, "")
	}

	// session
	cache, _ := c.Get("cache")
	uploadSessionRaw, _ := cache.(*freecache.Cache).Get([]byte("callback_" + j.ID))
	var uploadSession serializer.UploadSession
	json.Unmarshal(uploadSessionRaw, &uploadSession)

	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	file := filesystem.FileStream{
		File:        c.Request.Body,
		MimeType:    c.Request.Header.Get("Content-Type"),
		Name:        uploadSession.Name,
		Size:        uploadSession.Size,
		VirtualPath: uploadSession.VirtualPath,
	}

	fs.Upload(file)
}

func CreateDownloadSession(c *gin.Context) {
	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// get object id
	objectID := c.GetUint("object_id")

	// get file by id
	file := model.GetFileByID()

	// session
	// store file model on cache
	sessionID := util.RandStringRunes(16)
	cache, _ := c.Get("cache")
	b, _ := json.Marshal(file)
	cache.(*freecache.Cache).Set([]byte("download_"+sessionID), b, 60)

	// get download address
	downloadURL := fs.GetDownloadURL(objectID, sessionID)

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
	json.Unmarshal(b, target)

	// prepare for download
	rsc := fs.OpenFile(&target)
	defer fs.CloseFile(rsc)

	// send file
	c.Header("Content-Disposition", "attachment; filename=\""+url.PathEscape(target.Name)+"\"")
	http.ServeContent(c.Writer, c.Request, target.Name, target.UpdatedAt, rsc)
}
