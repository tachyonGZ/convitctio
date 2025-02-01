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

	fs.CreateDirectory(param.Path)

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

	dir := fs.OpenDirectory(param.Path)
	cd, cf := fs.ReadDirectory(dir)

	type DirObj struct {
		Name string `json:"name"`
	}

	dirObjs := make([]DirObj, 0, len(cd))

	for _, subDir := range cd {
		dirObjs = append(dirObjs, DirObj{
			Name: subDir.Name,
		})
	}

	type FileObj struct {
		Name string `json:"name"`
	}

	fileObjs := make([]FileObj, 0, len(cd))

	for _, subFile := range cf {
		fileObjs = append(fileObjs, FileObj{
			Name: subFile.Name,
		})
	}

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
