package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

/*Createid 生成唯一uuid
 */
func Createid() string {
	u1 := uuid.Must(uuid.NewV4(), nil)
	id := u1.String()
	return id
}

func main3() {
	for i := 0; i < 20; i++ {
		fmt.Println(Createid())

	}
}
