package controller

import (
	"conviction/filesystem"
	"conviction/model"

	"github.com/gin-gonic/gin"
)

func CreateDirectory(c *gin.Context) {
	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// binding
	var param struct {
		Path string `uri:"path" json:"path" binding:"required,min=1,max=65535"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, "")
	}

	// create
	fs.CreateDirectory(param.Path)

	// response
	c.JSON(200, "")
}

func ListDirectory(c *gin.Context) {
	// create file system
	u, _ := c.Get("user")
	fs := filesystem.NewFileSystem(u.(*model.User))

	// binding
	var param struct {
		Path string `uri:"path" json:"path" binding:"required,min=1,max=65535"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(200, "")
	}

	// read
	dir := fs.OpenDirectory(param.Path)
	cd, cf := fs.ReadDirectory(dir)

	// convert directory model to json object
	type DirObj struct {
		Name string `json:"name"`
	}
	dirObjs := make([]DirObj, 0, len(cd))
	for _, subDir := range cd {
		dirObjs = append(dirObjs, DirObj{
			Name: subDir.Name,
		})
	}

	// convert file model to json object
	type FileObj struct {
		Name string `json:"name"`
	}
	fileObjs := make([]FileObj, 0, len(cd))
	for _, subFile := range cf {
		fileObjs = append(fileObjs, FileObj{
			Name: subFile.Name,
		})
	}

	// response
	c.JSON(
		200,
		struct {
			directories []DirObj
			files       []FileObj
		}{
			directories: dirObjs,
			files:       fileObjs,
		})
}
