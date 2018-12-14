package controllers

import (
	"fmt"
	. "myproject/models"

	"github.com/astaxie/beego/orm"
)

func IfRoleInUser(roleId int64, user *User) (b bool) {
	b = false
	o := orm.NewOrm()
	m2m := o.QueryM2M(user, "Role")
	if m2m.Exist(&Role{Id: roleId}) {
		b = true
		fmt.Println("Role Exist")
		return
	}
	return
}

func IfPermInUser(permId int64, user *User) (b bool) {
	b = false
	o := orm.NewOrm()
	m2m := o.QueryM2M(user, "Permission")
	if m2m.Exist(&Permission{Id: permId}) {
		b = true
		fmt.Println("Permission Exist")
		return
	}
	return
}

func IfUserInRole(userId int64, role *Role) (b bool) {
	b = false
	o := orm.NewOrm()
	m2m := o.QueryM2M(role, "User")
	if m2m.Exist(&User{Id: userId}) {
		b = true
		fmt.Println("Tag Exist")
		return
	}
	return
}

func IfPermInRole(permId int64, role *Role) (b bool) {
	b = false
	o := orm.NewOrm()

	m2m := o.QueryM2M(role, "Permission")
	if m2m.Exist(&Permission{Id: permId}) {
		b = true
		fmt.Println("Permission Exist")
		return
	}
	return
}

func IfUserInPermission(userId int64, permission *Permission) (b bool) {
	b = false
	o := orm.NewOrm()

	m2m := o.QueryM2M(permission, "User")
	if m2m.Exist(&User{Id: userId}) {
		b = true
		fmt.Println("Permission Exist")
		return
	}
	return
}

func IfRoleInPerm(roleId int64, permission *Permission) (b bool) {
	b = false
	o := orm.NewOrm()

	m2m := o.QueryM2M(permission, "Role")
	if m2m.Exist(&Role{Id: roleId}) {
		b = true
		fmt.Println("Role Exist")
		return
	}
	return
}
