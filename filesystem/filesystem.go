package filesystem

import (
	adapter "conviction/filesystem/driver"
	"conviction/filesystem/driver/local"
	"conviction/model"
	"conviction/util"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
)

type FileSystem struct {
	Owner   model.User
	Target  model.File
	Adapter adapter.IFileSystemAdapter
}

type GlobalFileSystem struct {
	Target  model.File
	Adapter adapter.IFileSystemAdapter
}

var fileSystemPool = sync.Pool{
	New: func() any {
		return &FileSystem{}
	},
}

var globalFileSystemPool = sync.Pool{
	New: func() any {
		return &GlobalFileSystem{}
	},
}

func GetFileSystem() *FileSystem {
	return fileSystemPool.Get().(*FileSystem)
}

func GetGlobalFileSystem() *GlobalFileSystem {
	return globalFileSystemPool.Get().(*GlobalFileSystem)
}

func (fs *FileSystem) DispatchAdapter() error {
	fs.Adapter = local.FileSystemAdapter{}
	return nil
}

func (fs *GlobalFileSystem) DispatchAdapter() error {
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

func NewGlobalFileSystem(owner_id string) *GlobalFileSystem {
	fs := GetGlobalFileSystem()
	fs.DispatchAdapter()
	return fs
}

func (fs *FileSystem) GrenateSavePath(fmd *FileHead, pFile *model.File) (save_path string, err error) {
	//strUserID := strconv.FormatUint(uint64(fs.Owner.Model.ID), 10)
	//strPath := fmd.VirtualPath
	//strName := fmd.Name
	//return path.Join("upload", strUserID, strPath, strName)
	pDir, e := model.FindUserDirectory(fs.Owner.UUID, pFile.DirectoryUUID)
	if e != nil {
		err = fmt.Errorf(e.Error())
		return
	}
	save_path = path.Join("upload", fs.Owner.UUID, pDir.UUID, pFile.UUID)
	return
}

func (fs *FileSystem) Upload(head *FileHead, body IFileBody, placehodler_id string) {

	pPlaceholder, e := model.FindUserFile(fs.Owner.UUID, placehodler_id)
	if e != nil {
		return
	}
	// implement
	fs.Adapter.Put(body, head.SavePath, head.GetSize())

	// placeholder to file
	pPlaceholder.PlaceholderToFile()
}

func (fs *GlobalFileSystem) Download(owner_id string, file_id string) (rsc util.ReadSeekCloser) {
	file, _ := model.FindUserFile(owner_id, file_id)
	rsc, _ = fs.Adapter.Get(file.Path)
	return
}

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

func (fs *FileSystem) CreatePlaceHolder(pHead *FileHead, dir_id string) (placeholder_id string, err error) {
	// check name conflict
	if model.IsSameNameFileExist(fs.Owner.UUID, dir_id, pHead.Name) {
		err = fmt.Errorf("same name exist")
		return
	}

	// record
	pHolder := &model.File{
		DirectoryUUID: dir_id,
		OwnerUUID:     fs.Owner.UUID,

		Name: pHead.Name,
		Path: pHead.SavePath,
		Size: pHead.Size,
	}
	e := pHolder.Create()
	if e != nil {
		err = fmt.Errorf("create placeholder error")
		return
	}
	placeholder_id = pHolder.UUID

	// grenate save path
	pHead.SavePath, _ = fs.GrenateSavePath(pHead, pHolder)

	// implement
	fs.Adapter.Put(io.NopCloser(strings.NewReader("")), pHead.SavePath, pHead.GetSize())

	return
}

func (fs *FileSystem) CreateDirectoryByPath(dirPath string) (dir_id string) {
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
			Name:       dirName,
			OwnerUUID:  fs.Owner.UUID,
			ParentUUID: &currentDir.UUID,
		}

		pChild.Create()

		currentDir = pChild
	}

	dir_id = currentDir.UUID
	return
}

func (fs *FileSystem) OpenDirectory(dirPath string) (dir_id string, exists bool, err error) {
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
			exists = false
			return
		}

		currentDir = childDir
	}
	exists = true
	dir_id = currentDir.UUID
	return
}

func (fs *FileSystem) ReadDirectory2(dir *model.Directory) ([]model.Directory, []model.File) {
	childDir, _ := dir.GetChildDirectory()
	childFile, _ := dir.GetChildFile()
	return childDir, childFile
}
