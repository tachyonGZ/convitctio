package filesystem

import (
	"conviction/model"
	"errors"
)

func (fs *FileSystem) DeleteFile(fileID string) (err error) {

	// get path of file to delete
	file, err := model.FindUserFile(fs.Owner.UUID, fileID)
	if err != nil {
		return
	}
	path := file.Path

	// remove record
	err = model.DeleteUserFile(fs.Owner.ID, fileID)
	if err != nil {
		return
	}

	// driver
	fs.Adapter.Delete(path)
	return
}

func (fs *FileSystem) CreateSharedFile(sourceID string) (sharedFileID string, err error) {

	// check is file exist
	exist, e := model.IsUserOwnFile(fs.Owner.UUID, sourceID)
	if e != nil {
		err = errors.New("find file fail")
		return
	}
	if !exist {
		err = errors.New("file not exist")
		return
	}

	// create shared file
	newSharedFile := model.SharedFile{
		CreatorUUID: fs.Owner.UUID,
		SourceUUID:  sourceID,
	}
	if e = newSharedFile.Create(); e != nil {
		err = errors.New("create shared file fail")
		return
	}

	sharedFileID = newSharedFile.UUID
	return
}

func (fs *FileSystem) DeleteSharedFile(sharedFileID string) (err error) {

	// delete shared file
	return
}

func (fs *FileSystem) GetSharedFile(shared_file_id string) (file_id string, err error) {

	return
}
