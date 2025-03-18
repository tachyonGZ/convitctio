package memocache

import (
	"conviction/serializer"
	"encoding/json"

	"github.com/coocood/freecache"
)

var cacheSize = 100 * 1024 * 1024
var pCache = freecache.NewCache(cacheSize)

func SetDownloadSession(key string, session *serializer.DownloadSession, expire int) error {
	sessionRaw, _ := json.Marshal(session)
	return pCache.Set([]byte("download_"+key), sessionRaw, expire)
}

func SetUploadSession(key string, session *serializer.UploadSession, expire int) error {
	sessionRaw, _ := json.Marshal(session)
	return pCache.Set([]byte("callback_"+key), sessionRaw, expire)
}

func GetDownloadSession(key string) (*serializer.DownloadSession, error) {
	sessionRaw, err := pCache.Get([]byte("download_" + key))
	session := serializer.DownloadSession{}
	json.Unmarshal(sessionRaw, &session)
	return &session, err
}

func GetUploadSession(key string) (*serializer.UploadSession, error) {
	sessionRaw, err := pCache.Get([]byte("callback_" + key))
	session := serializer.UploadSession{}
	json.Unmarshal(sessionRaw, &session)
	return &session, err
}

func DeleteUploadSession(key string) bool {
	return pCache.Del(append([]byte("callback_"), key...))
}
