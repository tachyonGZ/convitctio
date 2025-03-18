package filesystem

import "conviction/model"

func (fs *FileSystem) DeleteFile(fileID string) (err error) {

	// get path of file to delete
	file, err := model.FindUserFile(fs.Owner.ID, fileID)
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
