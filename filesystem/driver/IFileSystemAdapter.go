package adapter

import (
	"conviction/serializer"
	"conviction/util"
	"io"
)

type IFileSystemAdapter interface {
	// @description		upload	a file
	// @param dst		path to store file
	// @param size		file size
	Put(file io.ReadCloser, dst string, size uint64) error

	// @description			delete one or more delicated file
	// @param filesPath		path of files
	Delete(filesPath []string) ([]string, error)

	// @description			get file content
	// @param filePath		path of file
	Get(filePath string) (util.ReadSeekCloser, error)

	Source(sessionId string) string

	Token(uploadSession serializer.UploadSession) serializer.UploadCredential
}
