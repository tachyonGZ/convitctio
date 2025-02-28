package main

import (
	"conviction/controller"
	"conviction/db"
	"conviction/middleware"
	"conviction/model"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
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

	store := memstore.NewStore([]byte("secret"))
	v1.Use(sessions.Sessions("convictio", store))
	v1.Use(middleware.CurrentUser())

	user := v1.Group("user")
	{
		user.POST("session", controller.UserLogin)
		user.POST("", controller.UserRegister)
	}

	auth := v1.Group("")
	auth.Use(middleware.AuthRequired())
	file := auth.Group("file")
	{

		file.PUT("upload", controller.CreateUploadSession)
		file.POST("upload", controller.UploadBySession)
		file.PUT("download", controller.CreateDownloadSession)
		file.GET("download", controller.DownloadBySession)
		file.DELETE("", controller.Delete)
	}

	directory := auth.Group("directory")
	{
		// create a directory
		directory.POST("", controller.CreateDirectory)

		// list all contents of a directory
		directory.GET("", controller.ListDirectory)
	}

	return r
}
