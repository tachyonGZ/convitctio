package serializer

type DownloadSessionDestType int

const (
	PersonalFile DownloadSessionDestType = iota
	SharedFile
)

type DownloadSession struct {
	Key string // session key

	DestType DownloadSessionDestType // destination type
	DestID   string                  // ID of destination
}
