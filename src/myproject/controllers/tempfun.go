package controllers

import (
	"fmt"
	. "myproject/models"

	"github.com/astaxie/beego/orm"
)

func IfRoleInUser(roleId int64, user *User) (b bool) {
	b = false
	o := orm.NewOrm()
	var roles []*Role
	var role = new(Role)
	_, err := o.QueryTable(role).Filter("User__User__Id", user.Id).All(&roles)
	if err != nil {
		fmt.Println("get permission failed err ", err)
		return
	}

	for _, p := range roles {
		if p.Id == roleId {
			b = true
			return
		}
	}
	return
}

func IfPermInUser(permId int64, user *User) (b bool) {
	b = false
	o := orm.NewOrm()
	var permissions []*Permission
	var permission = new(Permission)
	_, err := o.QueryTable(permission).Filter("User__User__Id", user.Id).All(&permissions)
	if err != nil {
		fmt.Println("get permission failed err ", err)
		return
	}

	for _, p := range permissions {
		if p.Id == permId {
			b = true
			return
		}
	}
	return
}
