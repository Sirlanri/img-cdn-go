package serves

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	DirInit()
}

/* DirInit 初始化文件夹
如果不存在所需目录，创建该目录*/
func DirInit() {

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Println("获取当前目录失败", err.Error())
	}

	dirs := []string{"pics"}

	for _, dir := range dirs {
		//是否存在文件夹
		flag := false
		for _, filedir := range files {
			if dir == filedir.Name() {
				flag = true
				break
			}
		}
		//如果遍历完扫描的文件夹不存在目标目录
		if !flag {
			err := os.Mkdir(dir, os.ModeDir)
			if err != nil {
				log.Fatalln("创建目录失败", err.Error())
			} else {
				log.Fatalln(dir, "文件夹不存在，创建成功")
			}
		}
	}

}
