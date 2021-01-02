package handlers

import (
	"imgcdn/serves"
	"log"
	"path/filepath"

	"github.com/kataras/iris/v12"
)

/*ImgUpload 上传单张图片
上传form的key为 "img"
*/
func ImgUpload(con iris.Context) {
	//最大上传图片限制为20M
	con.SetMaxRequestBodySize(20 * iris.MB)
	_, header, err := con.FormFile("img")
	if err != nil {
		log.Fatalln("图片上传失败！", err.Error())
		con.StatusCode(201)
		return
	}
	imgname := serves.Createid() + header.Filename
	path := filepath.Join("./pics", imgname)
	_, err = con.SaveFormFile(header, path)
	if err != nil {
		log.Fatalln("图片存入磁盘失败！", err.Error())
	} else {
		log.Fatalln("图片写入磁盘成功", imgname)
	}
}
