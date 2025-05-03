package controller

import (
	"conviction/filesystem"

	"github.com/gin-gonic/gin"
)

func CreateDirectory(c *gin.Context) {
	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// binding
	var param struct {
		ParentID string `json:"parent_id" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
	}

	// create dir
	newDirID := fs.CreateDirectory(param.ParentID, param.Name)

	// response
	c.JSON(
		200,
		struct {
			DirectoryID string `json:"directory_id"`
		}{
			DirectoryID: newDirID,
		})
}

func DeleteDirectory(c *gin.Context) {
	// data binding
	var param struct {
		DirectoryID string `json:"directory_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(500, err.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// delete directory
	if err := fs.DeleteDirectory(param.DirectoryID); err != nil {
		c.String(500, err.Error())
		return
	}

	// response
	c.String(200, "")
}

func GetDirectoryInfo(c *gin.Context) {
	// data binding
	var param struct {
		DirectoryID string `json:"directory_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.String(500, err.Error())
		return
	}

	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	// get directory info
	dirHead := fs.GetDirectoryHead(param.DirectoryID)

	// response
	c.JSON(
		200,
		struct {
			Name string `json:"name"`
			//OwnerName string `json:"owner_name"`
		}{
			Name: dirHead.Name,
		})
}

func ReadDirectory(c *gin.Context) {
	// create file system
	user_id, _ := c.Get("user_id")
	fs, _ := filesystem.NewFileSystem(user_id.(string))

	//var param struct {
	//	Path string `uri:"path" json:"path" binding:"required,min=1,max=65535"`
	//}
	// data binding
	var param struct {
		DirectoryID string `json:"directory_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(500, err.Error())
	}

	// read directory
	childDirID, childFileID := fs.ReadDirectory(param.DirectoryID)

	//	// convert directory model to json object
	//	type DirObj struct {
	//		Name string `json:"name"`
	//	}
	//	dirObjs := make([]DirObj, 0, len(cd))
	//	for _, subDir := range cd {
	//		dirObjs = append(dirObjs, DirObj{
	//			Name: subDir.Name,
	//		})
	//	}
	//
	//	// convert file model to json object
	//	type FileObj struct {
	//		Name string `json:"name"`
	//	}
	//	fileObjs := make([]FileObj, 0, len(cd))
	//	for _, subFile := range cf {
	//		fileObjs = append(fileObjs, FileObj{
	//			Name: subFile.Name,
	//		})
	//	}

	// response
	c.JSON(
		200,
		struct {
			Directories []string `json:"directories"`
			Files       []string `json:"files"`
		}{
			Directories: childDirID,
			Files:       childFileID,
		})
}
