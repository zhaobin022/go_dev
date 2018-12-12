package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

func init() {
	// 需要在init中注册定义的model
	fmt.Println("init models")
	orm.Debug = true

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(User), new(Role), new(Permission))

	// orm.RegisterDataBase("default", "mysql", "root:root@tcp(10.12.9.195:3306)/test?charset=utf8")
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")

	orm.RunSyncdb("default", false, true)

}
