package memocache

import (
	"encoding/json"

	"github.com/coocood/freecache"
)

var cacheSize = 100 * 1024 * 1024
var pCache = freecache.NewCache(cacheSize)

func SetJSON(key []byte, value any, expire int) error {
	valueM, _ := json.Marshal(value)
	return pCache.Set(key, valueM, expire)
}

func GetJSON(key []byte) (any, error) {
	valueRaw, err := pCache.Get(key)
	var value any
	json.Unmarshal(valueRaw, &value)
	return value, err
}
