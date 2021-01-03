package main

import (
	"imgcdn/handlers"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

func main() {
	app := iris.New()
	app.OnErrorCode(iris.StatusNotFound, handlers.NotFound)
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"https://*.ri-co.cn", "http://localhost:8080"}, //允许通过的主机名称
		//AllowedHeaders:   []string{"accept", "content-type", "authorization"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowCredentials: true,
		Debug:            true,
	})
	app.WrapRouter(crs.ServeHTTP)
	img := app.Party("/img")

	img.Post("/upload", handlers.ImgUploadOSS)
	img.Get("/del", handlers.DelImgOSS)

	img.Get("/test", handlers.Test)
	app.Run(iris.Addr(":8091"))

	return
}
