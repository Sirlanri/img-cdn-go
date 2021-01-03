package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db     *sql.DB
	dbpics []string

	bucket *oss.Bucket
)

func main() {
	//初始化数据库
	Db = ConnectDB()
	//获取picture的图片id
	getAllAddress()
	getusergroupURL()

	rmLocalPic()
	ossInit()
	uploadOSS()
}

//ConnectDB 初始化时，连接数据库
func ConnectDB() *sql.DB {
	Db, err := sql.Open("mysql", "root:123456@/whisper")
	if err != nil {
		fmt.Println("数据库初始化链接失败", err.Error())
	}

	if Db.Ping() != nil {
		fmt.Println("初始化-数据库-用户/密码/库验证失败", Db.Ping().Error())
		return nil
	}
	return Db
}
func getAllAddress() []string {
	tx, err := Db.Begin()
	if err != nil {
		fmt.Println("初始化", err.Error())
	}
	picsRows, err := tx.Query("select picaddress from picture")
	if err != nil {
		fmt.Println("数据库查询出错", err.Error())
		return nil
	}
	for picsRows.Next() {
		var pic string
		picsRows.Scan(&pic)
		pic = pic[25:]
		dbpics = append(dbpics, pic)
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("执行查询数据库图片出错", err.Error())
	}
	return dbpics
}
func getusergroupURL() {
	tx, _ := Db.Begin()
	//查询群组图片
	picsRows, err := tx.Query("select banner from igroup")
	if err != nil {
		fmt.Println("数据库查询出错", err.Error())
		return
	}
	for picsRows.Next() {
		var pic string
		picsRows.Scan(&pic)
		pic = pic[25:]
		fmt.Println(pic)
		dbpics = append(dbpics, pic)
	}

	//查询用户图片
	picsRows, err = tx.Query("select avatar,bannar from user")
	if err != nil {
		fmt.Println("数据库查询出错", err.Error())
		return
	}
	for picsRows.Next() {
		var pic, ban string
		picsRows.Scan(&pic, &ban)
		if pic != "" {
			pic = pic[25:]
			dbpics = append(dbpics, pic)

		}
		if ban != "" {
			ban = ban[25:]
			dbpics = append(dbpics, ban)
		}
	}
	tx.Commit()
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
func uploadOSS() {
	localpics, err := ioutil.ReadDir("./uploadpics")
	if err != nil {
		fmt.Println("OSS 读取本地失败", err.Error())
		return
	}
	for _, pic := range localpics {
		full := "./uploadpics/" + pic.Name()
		fmt.Println(full)
		err = bucket.PutObjectFromFile("img/"+pic.Name(), full)
		if err != nil {
			fmt.Println("上传oss出错", err.Error())

		}
	}
}
func rmLocalPic() {
	fmt.Println("开始删除操作")
	list, err := ioutil.ReadDir("./uploadpics")
	if err != nil {
		fmt.Println("读取失败", err.Error())
		return
	}
	var delids []string
	for _, v := range list {
		name := v.Name()
		flag := false
		for _, dbpic := range dbpics {
			if dbpic == name {
				//存在数据库，取消
				flag = true
				break
			}
		}
		if !flag {
			delids = append(delids, name)
		}
	}

	fmt.Println("磁盘中多余的图：", delids)

	//开始删除操作
	for _, v := range delids {
		err = os.Remove("./uploadpics/" + v)
		if err != nil {
			fmt.Println("删除失败", err.Error())
		}

	}
}

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
