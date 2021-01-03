package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

/*
type info struct {
	Endpoint   string `json:"endpoint"`
	Accessid   string `json:"accessid"`
	Secret     string `json:"secret"`
	Bucketname string `json:"bucketname"`
}
type jsonfile struct {
	Ossinfo info `json:"ossinfo"`
}
*/
func getinfo() {
	data, err := ioutil.ReadFile("../configs/secret.json")
	if err != nil {
		fmt.Println("读取json文件出错", err.Error())
		return
	}
	var totalinfo jsonfile

	err = json.Unmarshal(data, &totalinfo)
	if err != nil {
		fmt.Println("解析json出错", err.Error())
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

	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			fmt.Println("列举文件出错", err.Error())
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			fmt.Println("Bucket: ", object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	//上传文件

}

func main2() {
	getinfo()
}
