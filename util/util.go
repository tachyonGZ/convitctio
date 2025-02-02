package util

import (
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

func IsNotExist(name string) bool {
	_, err := os.Stat(name)

	if err == nil {
		return false
	}

	return os.IsNotExist(err)
}

// RelativePath 获取相对可执行文件的路径
func RelativePath(name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	e, _ := os.Executable()
	return filepath.Join(filepath.Dir(e), name)
}

// FormSlash 将path中的反斜杠'\'替换为'/'
func FormSlash(old string) string {
	return path.Clean(strings.ReplaceAll(old, "\\", "/"))
}

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// SplitPath 分割路径为列表
func SplitPath(fullPath string) []string {
	fullPath = path.Clean(fullPath)
	var MakeList func(*string) *[]string
	MakeList = func(p *string) *[]string {
		if *p == "/" {
			return new([]string)
		}

		base := path.Base(*p)
		*p = path.Dir(*p)

		l := MakeList(p)
		*l = append(*l, base)
		return l
	}

	return *MakeList(&fullPath)
}
