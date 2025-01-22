package filesystem

import (
	"io"
	"time"
)

type IFileStream interface {
	io.ReadCloser
	io.Seeker
	GetName() string
	GetSize() uint64
}

type FileStream struct {
	File   io.ReadCloser
	Seeker io.Seeker

	LastModified *time.Time
	MimeType     string
	Name         string
	SavePath     string
	Size         uint64
	VirtualPath  string
}

func (file FileStream) Read(p []byte) (n int, err error) {
	return file.File.Read(p)
}

func (file FileStream) Close() error {
	return file.File.Close()
}

func (file FileStream) GetName() string {
	return file.Name
}

func (file FileStream) GetSize() uint64 {
	return file.Size
}
