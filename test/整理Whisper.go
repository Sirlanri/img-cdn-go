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

//Db 创建的唯一指针
var (
	Db           *sql.DB
	dbpics       []string
	avatars      []string
	avatarssrc   []string
	bannars      []string
	bannarssrc   []string
	groupBans    []string
	groupBanssrc []string
	bucket       *oss.Bucket
)

func main4() {
	Db = ConnectDB()
	getAllAddress()
	getLocalPic()
	//getAllPostidFromPic()
	//dropTpoic()

	//ossInit()
	//uploadOSS()
	//changeURL()

	//changeusergroup()
	getusergroupURL()
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

//删除pic无关的id
func getAllPostidFromPic() {
	tx, _ := Db.Begin()
	var (
		picids  []string
		postids []string
		delids  []string
	)

	postrow, _ := tx.Query("select postid from post")
	for postrow.Next() {
		var postid string
		postrow.Scan(&postid)
		postids = append(postids, postid)
	}

	picrow, _ := tx.Query("select postid from picture")
	for picrow.Next() {
		var picid string
		picrow.Scan(&picid)
		flag := false
		//获取无效的id
		for _, id := range postids {
			if picid == id {
				flag = true
				break
			}
		}
		if !flag {
			delids = append(delids, picid)
		}
		picids = append(picids, picid)
	}
	/*
		for _, id := range delids {
			tx.Exec("delete from picture where postid=?", id)
		}
	*/
	tx.Commit()
	fmt.Println(delids)
}

func dropTpoic() {
	tx, _ := Db.Begin()
	var (
		postids []string
		delids  []string
	)

	postrow, _ := tx.Query("select postid from post")
	for postrow.Next() {
		var postid string
		postrow.Scan(&postid)
		postids = append(postids, postid)
	}
	picrow, _ := tx.Query("select postid from tag")
	for picrow.Next() {
		var topicid string
		picrow.Scan(&topicid)
		flag := false
		//获取无效的id
		for _, id := range postids {
			if topicid == id {
				flag = true
				break
			}
		}
		if !flag {
			delids = append(delids, topicid)
		}
	}

	for _, id := range delids {
		tx.Exec("delete from tag where postid=?", id)
	}

	tx.Commit()
	//fmt.Println(delids)
}

func dropReply() {
	tx, _ := Db.Begin()
	var (
		postids []string
		delids  []string
	)

	postrow, _ := tx.Query("select postid from post")
	for postrow.Next() {
		var postid string
		postrow.Scan(&postid)
		postids = append(postids, postid)
	}

	//postid
	picrow, _ := tx.Query("select postid from reply")
	for picrow.Next() {
		var topicid string
		picrow.Scan(&topicid)
		flag := false
		//获取无效的id
		for _, id := range postids {
			if topicid == id {
				flag = true
				break
			}
		}
		if !flag {
			delids = append(delids, topicid)
		}
	}
	for _, id := range delids {
		tx.Exec("delete from reply where postid=?", id)
	}

	//用户id
	var (
		valids []string
	)
	validRow, _ := tx.Query("select userid from user")
	for validRow.Next() {
		var valid string
		validRow.Scan(&valid)
		valids = append(valids, valid)
	}

	replyuseridRow, _ := tx.Query("select replyid, fromUser, toUser from reply")
	for replyuseridRow.Next() {
		var (
			user1, user2, replyid string
		)
		replyuseridRow.Scan(&replyid, &user1, &user2)

	}

	tx.Commit()
}

func getLocalPic() {
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

func changeURL() {
	fmt.Println("开始修改数据库图片地址")
	tx, _ := Db.Begin()
	urlsRow, err := tx.Query("select postid, picaddress from picture")
	if err != nil {
		fmt.Println("修改url 数据库查询出错", err.Error())
		return
	}
	var (
		sourceurls []string
		afterurls  []string
		ids        []string
	)
	for urlsRow.Next() {
		var url, id string
		urlsRow.Scan(&id, &url)
		sourceurls = append(sourceurls, url)
		afterurl := "https://cdn.ri-co.cn/img/" + url[37:]
		afterurls = append(afterurls, afterurl)
		ids = append(ids, id)
	}

	for index, url := range afterurls {
		_, err = tx.Exec("update picture set picaddress=? where picaddress=?", url, sourceurls[index])
		//_, err = tx.Exec("update user set avatar=? where avatar=?", url, sourceurls[index])
		if err != nil {
			fmt.Println("picture写入修改后的数据出错", err.Error())
		}
	}

	tx.Commit()
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

func changeusergroup() {
	tx, _ := Db.Begin()
	for index, url := range avatars {
		_, err := tx.Exec("update user set avatar=? where avatar=?", url, avatarssrc[index])
		if err != nil {
			fmt.Println("写入修改后的数据出错", err.Error())
		}
	}
	for index, url := range bannars {
		_, err := tx.Exec("update user set bannar=? where bannar=?", url, bannarssrc[index])
		if err != nil {
			fmt.Println("写入修改后的数据出错", err.Error())
		}
	}
	for index, url := range groupBans {
		_, err := tx.Exec("update igroup set bannar=? where bannar=?", url, groupBanssrc[index])
		if err != nil {
			fmt.Println("写入修改后的数据出错", err.Error())
		}
	}
	tx.Commit()
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
