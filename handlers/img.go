package handlers

import (
	"fmt"
	"imgcdn/serves"
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
		fmt.Println("图片上传失败！", err.Error())
		con.StatusCode(201)
		return
	}
	filename := header.Filename
	if len(filename) > 20 {
		filename = filename[len(filename)-20:]
	}
	imgname := serves.Createid() + filename
	path := filepath.Join("./pics", imgname)
	_, err = con.SaveFormFile(header, path)
	if err != nil {
		fmt.Println("图片存入磁盘失败！", err.Error())
	} else {
		fmt.Println("图片写入磁盘成功", imgname)
		fullurl := "https://api.ri-co.cn/cdn/getimg/" + imgname
		con.WriteString(fullurl)
	}
}
