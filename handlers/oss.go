package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/kataras/iris/v12"
)

var (
	bucket    *oss.Bucket
	accessKey string
	varify    = false
)

func init() {
	ossInit()
}

func ossInit() {
	//读取保密key内容
	data, err := ioutil.ReadFile("./secret.json")
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
	accessKey = totalinfo.Access
	point := totalinfo.Ossinfo.Endpoint
	keyid := totalinfo.Ossinfo.Accessid
	keysecret := totalinfo.Ossinfo.Secret
	bucketName := totalinfo.Ossinfo.Bucketname

	ossClient, err := oss.New(point, keyid, keysecret)
	if err != nil {
		fmt.Println("连接oss失败！")
		return
	}
	bucket, err = ossClient.Bucket(bucketName)
	if err != nil {
		fmt.Println("oss-连接bucket失败", err.Error())
		return
	}
}

/*ImgUploadOSS 上传单张图片到OSS
上传form的key为 "img"
*/
func ImgUploadOSS(con iris.Context) {
	//最大上传图片限制为20M
	con.SetMaxRequestBodySize(20 * iris.MB)
	if varify {
		//权限认证
		power := con.URLParamDefault("access", "")
		if power != accessKey {
			fmt.Println("无上传权限的请求，已拒绝")
			con.StatusCode(201)
			con.WriteString("你没权限哒~ 别折腾了，小可爱")
			return
		}
	}

	imgfile, header, err := con.FormFile("img")
	if err != nil {
		fmt.Println("图片上传失败！", err.Error())
		con.StatusCode(201)
		return
	}
	filename := GetFileName(header.Filename)

	err = bucket.PutObject("img/"+filename, imgfile)
	if err != nil {
		fmt.Println("上传文件至oss失败", err.Error())
		return
	}
	headurl := "https://cdn.ri-co.cn/img/"
	con.WriteString(headurl + filename)
}

/*DelImgOSS 删除OSS的图片
有安全问题，暂停使用此API
*/
func DelImgOSS(con iris.Context) {
	picurl := con.URLParamDefault("addr", "")
	if picurl == "" {
		con.StatusCode(201)
		con.WriteString("删除地址不存在")
		return
	}

	//权限认证
	if varify {
		power := con.URLParamDefault("access", "")
		if power == "" || power != accessKey {
			fmt.Println("无删除权限的请求，已拒绝")
			con.StatusCode(201)
			con.WriteString("你没权限哒~ 别折腾了，小可爱")
			return
		}
	}

	//剔除url前缀，得到bucket内的路径
	realuri := picurl[21:]
	fmt.Println(realuri)

	//删除图片
	err := bucket.DeleteObject(realuri)
	if err != nil {
		fmt.Println("OSS-删除失败", err.Error())
		con.StatusCode(202)
		con.WriteString("oss删除失败")
		return
	}
	con.WriteString("删除成功")

}
