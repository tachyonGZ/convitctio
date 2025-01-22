package filesystem

import (
	adapter "conviction/filesystem/driver"
	"conviction/model"
	"conviction/util"
	"path"
	"strconv"
	"sync"
)

var fileSystemPool = sync.Pool{
	New: func() any {
		return &FileSystem{}
	},
}

func GetFileSystem() *FileSystem {
	return fileSystemPool.Get().(*FileSystem)
}

type FileSystem struct {
	Owner   model.User
	Target  model.File
	Adapter adapter.IFileSystemAdapter
}

func NewFileSystem(Owner *model.User) *FileSystem {
	fs := GetFileSystem()
	fs.Owner = *Owner
	return fs
}

func (fs *FileSystem) GrenateSavePath(file FileStream) string {
	strUserID := strconv.FormatUint(uint64(fs.Owner.Model.ID), 10)
	strPath := file.VirtualPath
	strName := file.Name
	return path.Join("upload", strUserID, strPath, strName)
}

func (fs *FileSystem) Upload(file FileStream) {

	// grenate save path
	savePath := fs.GrenateSavePath(file)

	// implement
	fs.Adapter.Put(file, savePath, file.GetSize())
	fs.AfterUpload(file)
}

func (fs *FileSystem) AfterUpload(file FileStream) {
	target := model.File{
		UserID: fs.Owner.ID,
		Name:   file.Name,
		Path:   file.SavePath,
	}

	target.Insert()
}

func (fs *FileSystem) GetDownloadURL(objectID uint, sessionID string) string {

	source := fs.Adapter.Source(sessionID)

	return source
}

func (fs *FileSystem) OpenFile(target *model.File) util.ReadSeekCloser {

	rsc, _ := fs.Adapter.Get(target.Path)

	return rsc
}

func (fs *FileSystem) CloseFile(rsc util.ReadSeekCloser) {
	rsc.Close()
}
