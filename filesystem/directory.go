package filesystem

import (
	"conviction/model"
	"strconv"
)

type DirectoryHead struct {
	Name string
}

func (fs *FileSystem) CreateDirectory(parentID string, name string) (dirID string) {

	id64, _ := strconv.ParseUint(parentID, 10, 32)
	id32 := uint(id64)
	dir := &model.Directory{
		Name:     name + "/",
		OwnerID:  fs.Owner.ID,
		ParentID: &id32,
	}

	_ = dir.Create()

	dirID = strconv.FormatUint(uint64(dir.ID), 10)

	return
}

func (fs *FileSystem) DeleteDirectory(dirID string) (err error) {
	err = model.DeleteUserDirectory(fs.Owner.ID, dirID)
	return
}

func (fs *FileSystem) GetDirectoryHead(dirID string) (dirHead *DirectoryHead) {

	id64, _ := strconv.ParseUint(dirID, 10, 32)
	id32 := uint(id64)

	dir, _ := model.GetUserDirectory(fs.Owner.ID, id32)

	dirHead = &DirectoryHead{
		Name: dir.Name,
	}

	return
}

func (fs *FileSystem) ReadDirectory(dirID string) (childDirID []string, childFileID []string) {

	// get dir will be read
	dirID64, _ := strconv.ParseUint(dirID, 10, 32)
	dir, _ := model.GetUserDirectory(fs.Owner.ID, uint(dirID64))

	// get child dir
	childDirs, _ := dir.GetChildDirectory()
	for _, childDir := range childDirs {
		ID := strconv.FormatUint(uint64(childDir.ID), 10)
		childDirID = append(childDirID, ID)
	}

	// get child file
	childFiles, _ := dir.GetChildFile()
	for _, childFile := range childFiles {
		ID := strconv.FormatUint(uint64(childFile.ID), 10)
		childFileID = append(childFileID, ID)
	}
	return
}
