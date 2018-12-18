package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

var regStruct map[string]interface{}

func GetObjFromStr(name string) (obj interface{}, err error) {
	obj, ok := regStruct[name]
	if ok == true {
		return
	} else {
		err = fmt.Errorf("not find struct")
		return
	}
}

func init() {
	// 需要在init中注册定义的model
	fmt.Println("init models")
	orm.Debug = false

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(User), new(Role), new(Permission))

	// orm.RegisterDataBase("default", "mysql", "root:root@tcp(10.12.9.195:3306)/test?charset=utf8")
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(localhost:3306)/test?charset=utf8")

	orm.RunSyncdb("default", false, true)
	regStruct = make(map[string]interface{})
	regStruct["Permission"] = Permission{}
	regStruct["Role"] = Role{}
	regStruct["User"] = User{}

	PermUrlType[0] = "函数名"
	PermUrlType[1] = "Url地址"
	PermUrlType[2] = "正则表达式"
}
