package main

import (
	"conviction/config"
	"conviction/controller"
	"conviction/db"
	"conviction/middleware"
	"conviction/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	fmt.Println("helloworld")
	db.InitDB()
	model.Migration(db.GetDB())
}

func main() {

	// release database
	defer func() {
		db.ReleaseDB()
	}()

	// http server
	api := InitRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: api,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server...")
}

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api")

	// user session
	store := memstore.NewStore([]byte("secret"))
	v1.Use(sessions.Sessions("convictio", store))

	// disable cache
	v1.Use(middleware.NoCache())

	// rate limiter
	v1.Use(middleware.RateLimit(10, 1))

	user := v1.Group("user")
	{
		user.POST("login", controller.UserLogin)
		user.POST("register", controller.UserRegister)
	}

	auth := v1.Group("")
	auth.Use(middleware.AuthRequired())
	{
		file := auth.Group("file")
		{

			file.POST("upload/session", controller.CreateUploadSession)
			file.POST("download/session", controller.CreateDownloadSession)

			// delete a file
			file.POST("delete", controller.DeleteFile)

			// get info of file
			file.POST("info", controller.GetFileStatus)

			// move a file
			file.POST("move", controller.MoveFile)

			// rename a file
			file.POST("rename", controller.RenameFile)
		}

		session := auth.Group("session")
		{
			session.GET("download/:session_id", controller.Download)
			session.POST("upload/:session_id", controller.Upload)
		}

		directory := auth.Group("directory")
		{
			// create a directory
			directory.POST("create", controller.CreateDirectory)

			// delete a directory
			directory.POST("delete", controller.DeleteDirectory)

			// get info of directory
			directory.POST("info", controller.GetDirectoryInfo)

			// read content of directory
			directory.POST("read", controller.ReadDirectory)

		}

		share := auth.Group("share")
		{
			share.POST("create", controller.CreateSharedFile)
			share.POST("delete", controller.DeleteSharedFile)
		}
	}
	return r
}
