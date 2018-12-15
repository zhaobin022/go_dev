package controllers

import (
	"fmt"
	. "myproject/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func a() {
	fmt.Println("init controller")
	o := orm.NewOrm()
	var permission1 = &Permission{Name: "create_user"}
	if created, id, err := o.ReadOrCreate(permission1, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		} else {
			fmt.Println("Get an object. Id:", id)
		}
	}

	var permission2 = &Permission{Name: "update_user"}
	if created, id, err := o.ReadOrCreate(permission2, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		} else {
			fmt.Println("Get an object. Id:", id)
		}
	}

	encryPassFormt := GetEncryPass("root")
	var user1 = &User{Name: "user1"}

	user1.Password = encryPassFormt

	if created, id, err := o.ReadOrCreate(user1, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		} else {
			fmt.Println("Get an object. Id:", id)
		}
	}

	var user2 = &User{Name: "user2"}
	user2.IsAdmin = true
	if created, id, err := o.ReadOrCreate(user2, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		} else {
			fmt.Println("Get an object. Id:", id)
		}
	}

	m2m := o.QueryM2M(user1, "Permission")
	num, err := m2m.Add(permission1)
	if err == nil {
		fmt.Println("Added nums: ", num)
	}

	var role1 = &Role{Name: "role1"}
	if created, id, err := o.ReadOrCreate(role1, "Name"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		} else {
			fmt.Println("Get an object. Id:", id)
		}
	}

	m2m = o.QueryM2M(role1, "Permission")
	num, err = m2m.Add(permission1)
	if err == nil {
		fmt.Println("Added nums: ", num)
	}

	m2m = o.QueryM2M(role1, "User")
	num, err = m2m.Add(user2)
	if err == nil {
		fmt.Println("Added nums: ", num)
	}

	// _ := user2.IfhasPermisson("create_user")

}

type Response struct {
	Status bool
	Msg    string
}

type BaseControl struct {
	beego.Controller
}

func (this *BaseControl) Prepare() {

	var username = this.GetSession("username")
	if username != "" {
		this.Data["username"] = username
	}
}

func init() {
	beego.AddFuncMap("IfObjInObjRel", IfObjInObjRel)
}
