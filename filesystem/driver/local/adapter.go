package local

import (
	"conviction/serializer"
	"conviction/util"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FileSystemAdapter struct {
}

func (fsa FileSystemAdapter) Put(srcFile io.ReadCloser, dst string, size uint64) error {
	defer srcFile.Close()
	dst = util.RelativePath(dst)

	// create dir
	dstPath := filepath.Dir(dst)

	if util.IsNotExist(dstPath) {
		err := os.MkdirAll(dstPath, 0744)
		if err != nil {
			return err
		}
	}

	// create file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// write content
	_, err = io.Copy(dstFile, srcFile)

	return err
}

func (fsa FileSystemAdapter) Delete(filePath []string) ([]string, error) {

	failedPath := make([]string, 0, len(filePath))
	var retErr error

	for _, path := range filePath {
		err := os.Remove(util.RelativePath(filepath.FromSlash(path)))
		if err != nil {
			retErr = err
			failedPath = append(failedPath, path)
		}
	}

	return failedPath, retErr
}

func (fsa FileSystemAdapter) Get(filePath string) (util.ReadSeekCloser, error) {
	file, err := os.Open(util.RelativePath(filePath))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (fsa FileSystemAdapter) Source(sessionID string) string {

	URI, _ := url.Parse(fmt.Sprintf("/api/v3/file/download/%s", sessionID))

	return URI.String()
}

// func (fsa FileSystemAdapter) Token(uploadSession *serializer.UploadSession) *serializer.UploadCredential {
func (fsa FileSystemAdapter) Token(uploadSession *serializer.UploadSession) *serializer.UploadCredential {
	if !util.IsNotExist(uploadSession.SavePath) {
		return nil
	}

	return &serializer.UploadCredential{
		SessionID: uploadSession.Key,
	}
}

func (fsa FileSystemAdapter) IsFileExist(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return false
	}

	return !os.IsNotExist(err)
}
