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
	if len(filename) > 30 {
		filename = filename[len(filename)-30:]
	}
	imgname := serves.Createid() + filename
	path := filepath.Join("./pics", imgname)
	_, err = con.SaveFormFile(header, path)
	if err != nil {
		fmt.Println("图片存入磁盘失败！", err.Error())
	} else {
		fmt.Println("图片写入磁盘成功", imgname)
		fullurl := "https://cdn.ri-co.cn/img/get/" + imgname
		con.WriteString(fullurl)
	}
}

/*GetFileName 为上传文件设置唯一文件名
生成uuid+截取后的文件名*/
func GetFileName(source string) string {
	if len(source) > 30 {
		source = source[len(source)-30:]
	}
	uuid := serves.Createid()
	return uuid + source
}

//结构图区域

type info struct {
	Endpoint   string `json:"endpoint"`
	Accessid   string `json:"accessid"`
	Secret     string `json:"secret"`
	Bucketname string `json:"bucketname"`
}
type jsonfile struct {
	Ossinfo info   `json:"ossinfo"`
	Access  string `json:"access"`
}
