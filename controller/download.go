package controller

import (
	"conviction/filesystem"
	"conviction/memocache"
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
	fs := filesystem.NewGlobalFileSystem(user_id.(string))

	// get download session from cache
	session, err := memocache.GetDownloadSession(param.SessionID)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// prepare for download
	rsc := fs.Download(session.OwnerID, session.FileID)
	defer rsc.Close()

	// send file
	c.Header(
		"Content-Disposition",
		"attachment; filename=\""+url.PathEscape(session.Name)+"\"")
	http.ServeContent(c.Writer, c.Request, session.Name, time.Time{}, rsc)

	c.String(200, "")
}
