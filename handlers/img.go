package handlers

import (
	"encoding/json"
	"fmt"
	"imgcdn/serves"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
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
	imgfile, header, err := con.FormFile("img")
	if err != nil {
		fmt.Println("图片上传失败！", err.Error())
		con.StatusCode(201)
		return
	}
	filename := GetFileName(header.Filename)
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

	point := totalinfo.Ossinfo.Endpoint
	keyid := totalinfo.Ossinfo.Accessid
	keysecret := totalinfo.Ossinfo.Secret
	bucketName := totalinfo.Ossinfo.Bucketname

	ossClient, err := oss.New(point, keyid, keysecret)
	if err != nil {
		fmt.Println("连接oss失败！")
		return
	}
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		fmt.Println("oss-连接bucket失败", err.Error())
		return
	}
	myfile := io.MultiReader(imgfile)
	err = bucket.PutObject("imgtest/"+filename, myfile)
	if err != nil {
		fmt.Println("上传文件至oss失败", err.Error())
		return
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
	Ossinfo info `json:"ossinfo"`
}
