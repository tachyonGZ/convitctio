package main

import (
	"conviction/controller"
	"conviction/db"
	middlewware "conviction/middleware"
	"conviction/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
	db.InitDB()
	model.Migration(db.GetDB())
}

func main() {
	api := InitRouter()
	api.Run(":8080")
}

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api")

	store := memstore.NewStore([]byte("secret"))
	v1.Use(sessions.Sessions("convictio", store))
	v1.Use(middlewware.CurrentUser())
	v1.Use(middlewware.MemoCache())

	user := v1.Group("user")
	{
		user.POST("session", controller.UserLogin)
		user.POST("", controller.UserRegister)
	}

	auth := v1.Group("")
	auth.Use(middlewware.AuthRequired())
	file := auth.Group("file")
	{

		file.PUT("upload", controller.CreateUploadSession)
		file.POST("upload", controller.UploadBySession)
		file.PUT("download", controller.CreateDownloadSession)
		file.GET("download", controller.DownloadBySession)
	}

	return r
}
