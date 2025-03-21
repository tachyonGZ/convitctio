package filesystem

import (
	adapter "conviction/filesystem/driver"
	"conviction/filesystem/driver/local"
	"conviction/model"
	"conviction/util"
	"io"
	"path"
	"strconv"
	"strings"
	"sync"
)

type FileSystem struct {
	Owner   model.User
	Target  model.File
	Adapter adapter.IFileSystemAdapter
}

var fileSystemPool = sync.Pool{
	New: func() any {
		return &FileSystem{}
	},
}

func GetFileSystem() *FileSystem {
	return fileSystemPool.Get().(*FileSystem)
}

func (fs *FileSystem) DispatchAdapter() error {
	fs.Adapter = local.FileSystemAdapter{}
	return nil
}

func NewFileSystem(Owner *model.User) *FileSystem {
	fs := GetFileSystem()
	fs.Owner = *Owner

	fs.DispatchAdapter()
	return fs
}

func NewFileSystem2(owner_id string) (fs *FileSystem, err error) {
	fs = GetFileSystem()

	pOwner, e := model.FindUser(owner_id)
	if e != nil {
		err = e
		return
	}
	fs.Owner = *pOwner

	fs.DispatchAdapter()

	return
}

func (fs *FileSystem) GrenateSavePath(fmd *FileHead) string {
	strUserID := strconv.FormatUint(uint64(fs.Owner.Model.ID), 10)
	strPath := fmd.VirtualPath
	strName := fmd.Name
	return path.Join("upload", strUserID, strPath, strName)
}

func (fs *FileSystem) Upload(head *FileHead, body IFileBody, pPlaceHolder *model.File) {

	// grenate save path
	savePath := fs.GrenateSavePath(head)

	// implement
	fs.Adapter.Put(body, savePath, head.GetSize())

	// placeholder to file
	pPlaceHolder.PlaceholderToFile()
}

func (fs *FileSystem) Download(file_id string) (rsc util.ReadSeekCloser) {
	file, _ := model.FindUserFile2(fs.Owner.UUID, file_id)
	rsc, _ = fs.Adapter.Get(file.Path)
	return
}

//func (fs *FileSystem) AfterUpload(head *FileHead) {
//
//	// record
//	db.GetDB().Create(
//		&model.File{
//			UserID: fs.Owner.ID,
//			Name:   head.Name,
//			Path:   head.SavePath,
//		})
//}

func (fs *FileSystem) GetDownloadURL(objectID uint, sessionID string) string {

	source := fs.Adapter.Source(sessionID)

	return source
}

func (fs *FileSystem) GetFileHead(fileID uint) (head *FileHead) {

	file, _ := model.GetFileByID(fileID, fs.Owner.ID)
	head = &FileHead{
		Name: file.Name,
	}
	return
}

func (fs *FileSystem) UpdateFile(target *model.File, file FileStream) {
	realPath := target.Path

	fs.Adapter.Put(file, realPath, file.GetSize())
}

func (fs *FileSystem) CreatePlaceHolder(head *FileHead, dir *model.Directory) (*model.File, error) {
	// grenate save path
	head.SavePath = fs.GrenateSavePath(head)

	// implement
	fs.Adapter.Put(io.NopCloser(strings.NewReader("")), head.SavePath, head.GetSize())

	// record
	placeholder := model.File{}
	err := fs.RecordFile(head, dir, &placeholder)
	return &placeholder, err
}

func (fs *FileSystem) RecordFile(head *FileHead, dir *model.Directory, file *model.File) error {
	*file = model.File{
		Name:        head.Name,
		Path:        head.SavePath,
		UserID:      fs.Owner.ID,
		Size:        head.Size,
		DirectoryID: dir.ID,
	}

	err := file.Create()
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileSystem) CreateDirectoryByPath(dirPath string) *model.Directory {
	dirPath = path.Clean(dirPath)
	pathList := util.SplitPath(dirPath)

	var currentDir *model.Directory
	for _, dirName := range pathList {

		if dirName == "/" {
			currentDir, _ = fs.Owner.Root()
			continue
		}

		pChild, exist, _ := currentDir.GetChild(dirName)
		if exist {
			currentDir = pChild
			continue
		}

		pChild = &model.Directory{
			Name:     dirName,
			OwnerID:  fs.Owner.ID,
			ParentID: &currentDir.ID,
		}

		pChild.Create()

		currentDir = pChild
	}

	return currentDir
}

func (fs *FileSystem) OpenDirectory(dirPath string) (*model.Directory, bool, error) {
	dirPath = path.Clean(dirPath)
	pathList := util.SplitPath(dirPath)

	var currentDir *model.Directory
	for _, dirName := range pathList {

		if dirName == "/" {
			currentDir, _ = fs.Owner.Root()
			continue
		}

		childDir, exist, _ := currentDir.GetChild(dirName)

		if !exist {
			return nil, false, nil
		}

		currentDir = childDir
	}
	return currentDir, true, nil
}

func (fs *FileSystem) ReadDirectory2(dir *model.Directory) ([]model.Directory, []model.File) {
	childDir, _ := dir.GetChildDirectory()
	childFile, _ := dir.GetChildFile()
	return childDir, childFile
}

func (fs *FileSystem) IsSameNameFileExists(name string, dir *model.Directory) bool {
	return model.IsSameNameFileExist(name, dir.ID, dir.OwnerID)
}
