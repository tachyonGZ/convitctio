package serializer

// UploadCredential 返回给客户端的上传凭证
type UploadCredential struct {
	SessionID string `json:"sessionID"`
	//ChunkSize   uint64   `json:"chunkSize"` // 分块大小，0 为部分快
	Expires     int64    `json:"expires"` // 上传凭证过期时间， Unix 时间戳
	UploadURLs  []string `json:"uploadURLs,omitempty"`
	Credential  string   `json:"credential,omitempty"`
	UploadID    string   `json:"uploadID,omitempty"`
	Callback    string   `json:"callback,omitempty"` // 回调地址
	Path        string   `json:"path,omitempty"`     // 存储路径
	AccessKey   string   `json:"ak,omitempty"`
	KeyTime     string   `json:"keyTime,omitempty"` // COS用有效期
	Policy      string   `json:"policy,omitempty"`
	CompleteURL string   `json:"completeURL,omitempty"`
}
