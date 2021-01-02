package handlers

import (
	"encoding/json"
	"fmt"
	"imgcdn/serves"
	"io/ioutil"
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

func ImgUploadOSS(con iris.Context) {
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

	//读取保密key内容
	data, err := ioutil.ReadFile("./configs/secret.json")
	if err != nil {
		fmt.Println("读取json文件出错", err.Error())
		return
	}
	var totalinfo jsonfile

	err = json.Unmarshal(data, &totalinfo)
	if err != nil {
		fmt.Println("解析json出错", err.Error())
		return
	}

}

type info struct {
	Endpoint   string `json:"endpoint"`
	Accessid   string `json:"accessid"`
	Secret     string `json:"secret"`
	Bucketname string `json:"bucketname"`
}
type jsonfile struct {
	Ossinfo info `json:"ossinfo"`
}
