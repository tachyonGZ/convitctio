package filesystem

import (
	adapter "conviction/filesystem/driver"
	"conviction/filesystem/driver/local"
	"conviction/model"
	"conviction/util"
	"fmt"
	"sync"
)

type GuestFileSystem struct {
	Owner   *model.User
	Target  model.File
	Adapter adapter.IFileSystemAdapter
}

var globalFileSystemPool = sync.Pool{
	New: func() any {
		return &GuestFileSystem{}
	},
}

func GetGuestFileSystem() *GuestFileSystem {
	return globalFileSystemPool.Get().(*GuestFileSystem)
}

func (fs *GuestFileSystem) DispatchAdapter() error {
	fs.Adapter = local.FileSystemAdapter{}
	return nil
}

func NewGuestFileSystem(owner_id string) *GuestFileSystem {
	fs := GetGuestFileSystem()
	fs.DispatchAdapter()
	return fs
}

func (fs *GuestFileSystem) Download(shared_file_id string) (rsc util.ReadSeekCloser, err error) {
	p_shared_file, e := model.FindSharedFile(shared_file_id)
	if e != nil {
		err = fmt.Errorf("find shared file fail")
		return
	}

	file, e := model.FindUserFile(p_shared_file.CreatorUUID, p_shared_file.SourceUUID)
	if e != nil {
		err = fmt.Errorf("find user file fail")
		return
	}

	rsc, _ = fs.Adapter.Get(file.Path)
	return
}

func (fs *GuestFileSystem) GetFileHead(shared_file_id string) (p_head *FileHead, err error) {

	p_shared_file, e := model.FindSharedFile(shared_file_id)
	if e != nil {
		err = fmt.Errorf("invalid shared file id")
		return
	}

	file, e := model.FindUserFile(p_shared_file.SourceUUID, fs.Owner.UUID)
	if e != nil {
		err = fmt.Errorf("invalid file id")
	}

	p_head = &FileHead{
		Name: file.Name,
	}
	return
}
