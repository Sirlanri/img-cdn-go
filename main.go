package main

import (
	"imgcdn/handlers"
	"imgcdn/serves"

	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

func main() {
	app := iris.New()
	//初始化文件夹
	serves.DirInit()
	app.OnErrorCode(iris.StatusNotFound, handlers.NotFound)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowedHeaders:   []string{"accept", "content-type", "authorization"},
		AllowCredentials: false,
		Debug:            true,
	})
	app.WrapRouter(crs.ServeHTTP)
	img := app.Party("/img")

	img.Post("/upload", handlers.ImgUploadOSS)
	img.Get("/del", handlers.DelImgOSS)

	app.Run(iris.Addr(":8090"))

	return
}
