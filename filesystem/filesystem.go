package filesystem

import (
	adapter "conviction/filesystem/driver"
	"conviction/model"
	"conviction/util"
	"io"
	"path"
	"strconv"
	"strings"
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

func (fs *FileSystem) GrenateSavePath(fmd *FileHead) string {
	strUserID := strconv.FormatUint(uint64(fs.Owner.Model.ID), 10)
	strPath := fmd.VirtualPath
	strName := fmd.Name
	return path.Join("upload", strUserID, strPath, strName)
}

func (fs *FileSystem) Upload(head *FileHead, body IFileBody) {

	// grenate save path
	savePath := fs.GrenateSavePath(head)

	// implement
	fs.Adapter.Put(body, savePath, head.GetSize())

	fs.AfterUpload(head)
}

func (fs *FileSystem) AfterUpload(head *FileHead) {
	target := model.File{
		UserID: fs.Owner.ID,
		Name:   head.Name,
		Path:   head.SavePath,
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

func (fs *FileSystem) UpdateFile(target *model.File, file FileStream) {
	realPath := target.Path

	fs.Adapter.Put(file, realPath, file.GetSize())
}

func (fs *FileSystem) CreatePlaceHolder(fmd *FileHead) {
	// grenate save path
	fmd.SavePath = fs.GrenateSavePath(fmd)

	// implement
	fs.Adapter.Put(io.NopCloser(strings.NewReader("")), fmd.SavePath, fmd.GetSize())
}

func (fs *FileSystem) CreateDirectory(dirPath string) *model.Directory {

	dirPath = path.Clean(dirPath)
	pathList := util.SplitPath(dirPath)

	var currentDir *model.Directory
	for _, dirName := range pathList {

		var err error
		currentDir, err = currentDir.GetChild(dirName)
		if err != nil {
			continue
		}

		currentDir = &model.Directory{
			Name:     dirName,
			OwnerID:  fs.Owner.ID,
			ParentID: currentDir.ID,
		}

		currentDir.Create()
	}

	return currentDir
}
