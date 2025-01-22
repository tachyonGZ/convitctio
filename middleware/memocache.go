package middlewware

import (
	"github.com/coocood/freecache"
	"github.com/gin-gonic/gin"
)

var cacheSize = 100 * 1024 * 1024
var cache = freecache.NewCache(cacheSize)

func MemoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", &cache)
		c.Next()
	}
}
