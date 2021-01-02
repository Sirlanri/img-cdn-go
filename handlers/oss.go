package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/kataras/iris/v12"
)

/*ImgUploadOSS 上传单张图片到OSS
上传form的key为 "img"
*/
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
	err = bucket.PutObject("img/"+filename, myfile)
	if err != nil {
		fmt.Println("上传文件至oss失败", err.Error())
		return
	}
	headurl := "https://cdn.ri-co.cn/img/"
	con.WriteString(headurl + filename)
}
