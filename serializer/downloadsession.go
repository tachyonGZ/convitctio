package serializer

type DownloadSession struct {
	Key string // session key

	FileID  string // ID of dest file
	Name    string // 文件名
	OwnerID string
}
