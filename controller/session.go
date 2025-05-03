package controller

import (
	"conviction/filesystem"
	"conviction/memocache"
	"conviction/serializer"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

func Download(c *gin.Context) {

	// data binding
	var param struct {
		SessionID string `uri:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
		return
	}

	// create global file system
	user_id, _ := c.Get("user_id")

	// get download session from cache
	session, err := memocache.GetDownloadSession(param.SessionID)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// prepare for download
	if session.DestType == serializer.SharedFile {
		// create file system
		fs := filesystem.NewGuestFileSystem(user_id.(string))

		// get head
		p_head, e := fs.GetFileHead(session.DestID)
		if e != nil {
			c.String(500, e.Error())
			return
		}

		// get stream
		rsc, e := fs.Download(session.DestID)
		if e != nil {
			c.String(500, e.Error())
			return
		}
		defer rsc.Close()

		// send file
		c.Header(
			"Content-Disposition",
			"attachment; filename=\""+url.PathEscape(p_head.Name)+"\"")
		http.ServeContent(c.Writer, c.Request, "file", time.Time{}, rsc)
	} else if session.DestType == serializer.PersonalFile {
		// create file system
		fs, _ := filesystem.NewFileSystem(user_id.(string))

		// get head
		p_head, e := fs.GetFileHead(session.DestID)
		if e != nil {
			c.String(500, e.Error())
			return
		}

		// get stream
		rsc, e := fs.Download(session.DestID)
		if e != nil {
			c.String(500, e.Error())
			return
		}
		defer rsc.Close()

		// send file
		c.Header(
			"Content-Disposition",
			"attachment; filename=\""+url.PathEscape(p_head.Name)+"\"")
		http.ServeContent(c.Writer, c.Request, "file", time.Time{}, rsc)
	}

	c.String(200, "")
}

func Upload(c *gin.Context) {
	// data binding
	var param struct {
		SessionID string `uri:"session_id" binding:"required"`
	}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(500, err.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

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
