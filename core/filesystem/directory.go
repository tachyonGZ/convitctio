package filesystem

import (
	"conviction/model"
)

type DirectoryHead struct {
	Name string
}

func (fs *FileSystem) CreateDirectory(parentID string, name string) (dirID string) {

	dir := &model.Directory{
		Name:       name + "/",
		OwnerUUID:  fs.Owner.UUID,
		ParentUUID: &parentID,
	}

	_ = dir.Create()

	dirID = dir.UUID
	return
}

func (fs *FileSystem) DeleteDirectory(dirID string) (err error) {
	err = model.DeleteUserDirectory(fs.Owner.UUID, dirID)
	return
}

func (fs *FileSystem) GetDirectoryHead(dirID string) (dirHead *DirectoryHead) {

	dir, _ := model.FindUserDirectory(fs.Owner.UUID, dirID)

	dirHead = &DirectoryHead{
		Name: dir.Name,
	}

	return
}

func (fs *FileSystem) ReadDirectory(dirID string) (childDirID []string, childFileID []string) {

	// get dir will be read
	dir, _ := model.FindUserDirectory(fs.Owner.UUID, dirID)

	// get child dir
	childDirs, _ := dir.GetChildDirectory()
	for _, childDir := range childDirs {
		childDirID = append(childDirID, childDir.UUID)
	}

	// get child file
	childFiles, _ := dir.GetChildFile()
	for _, childFile := range childFiles {
		childFileID = append(childFileID, childFile.UUID)
	}
	return
}
