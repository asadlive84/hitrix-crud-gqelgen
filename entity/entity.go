// entity/user.go
package entity

import "git.ice.global/packages/beeorm/v4"

type User struct {
	beeorm.ORM
	ID   uint64 `orm:"localCache"`
	Name string `orm:"required;length=100"`
	Age  int    `orm:"required"`
}

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(&User{})
}
