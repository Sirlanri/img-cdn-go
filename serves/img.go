package serves

import uuid "github.com/satori/go.uuid"

/*Createid 生成唯一uuid
 */
func Createid() string {
	u1 := uuid.Must(uuid.NewV4(), nil)
	id := u1.String()
	return id
}
