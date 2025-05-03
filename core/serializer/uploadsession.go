package serializer

import "time"

type UploadSession struct {
	Key string // session key

	PlaceholderID  string // file id
	OwnerID        string // user id
	VirtualPath    string // 用户文件路径，不含文件名
	MimeType       string
	Name           string     // 文件名
	Size           uint64     // 文件大小
	SavePath       string     // 物理存储路径，包含物理文件名
	LastModified   *time.Time // 可选的文件最后修改日期
	Callback       string     // 回调 URL 地址
	CallbackSecret string     // 回调 URL
	UploadURL      string
	UploadID       string
	Credential     string
}
