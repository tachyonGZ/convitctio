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

func NewFileSystem(owner_id string) (fs *FileSystem, err error) {
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

func (fs *FileSystem) Download(file_id string) (rsc util.ReadSeekCloser, err error) {

	file, e := model.FindUserFile(fs.Owner.UUID, file_id)
	if e != nil {
		err = fmt.Errorf("find user file fail")
		return
	}

	rsc, _ = fs.Adapter.Get(file.Path)
	return
}

func (fs *FileSystem) GrenateSavePath(fmd *FileHead, pFile *model.File) (save_path string, err error) {
	//strUserID := strconv.FormatUint(uint64(fs.Owner.Model.ID), 10)
	//strPath := fmd.VirtualPath
	//strName := fmd.Name
	//return path.Join("upload", strUserID, strPath, strName)
	pDir, e := model.FindUserDirectory(fs.Owner.UUID, pFile.DirectoryUUID)
	if e != nil {
		err = fmt.Errorf("not find direcotory")
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

func (fs *FileSystem) GetDownloadURL(objectID uint, sessionID string) string {

	source := fs.Adapter.Source(sessionID)

	return source
}

func (fs *FileSystem) GetFileHead(file_id string) (p_head *FileHead, err error) {

	file, e := model.FindUserFile(file_id, fs.Owner.UUID)
	if e != nil {
		err = fmt.Errorf("invalid file id")
	}

	p_head = &FileHead{
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

func (fs *FileSystem) CreateDirectoryByPath(dirPath string) (dir_id string, err error) {
	dirPath = path.Clean(dirPath)
	pathList := util.SplitPath(dirPath)

	var currentDir *model.Directory
	for _, dirName := range pathList {

		if dirName == "/" {
			pRoot, e := fs.Owner.Root()
			if e != nil {
				err = fmt.Errorf("get root fail")
				return
			}

			currentDir = pRoot
			continue
		}

		pChild, exist, e := currentDir.GetChild(dirName)
		if e != nil {
			err = fmt.Errorf("get child fail")
			return
		}
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
			pRoot, e := fs.Owner.Root()
			if e != nil {
				err = fmt.Errorf("get root fail")
				return
			}

			currentDir = pRoot
			continue
		}

		childDir, exist, e := currentDir.GetChild(dirName)
		if e != nil {
			err = fmt.Errorf("get child fail")
			return
		}
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

func (fs *FileSystem) RenameFile(file_id string, name string) (err error) {
	pFile, e := model.FindUserFile(fs.Owner.UUID, file_id)
	if e != nil {
		err = fmt.Errorf("invalid file id")
		return
	}

	e = pFile.Rename(name)
	if e != nil {
		err = fmt.Errorf("file rename fail")
		return
	}

	return
}
