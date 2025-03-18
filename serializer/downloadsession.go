package serializer

type DownloadSession struct {
	FileID uint   // ID of dest file
	Key    string // UUID
	Name   string // 文件名
}
