package controllers

import (
	"fmt"
	. "myproject/models"
	"reflect"
	"strings"

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

func IfObjInObjRel(obj, objrel interface{}) (b bool) {
	b = false
	objType := reflect.TypeOf(obj)
	typeSlice := strings.Split(objType.String(), ".")
	typeName := typeSlice[len(typeSlice)-1]
	/*
		动态创建struct 很有用的功能
		t := reflect.ValueOf(regStruct[typeName]).Type()
		obj := reflect.New(t).Elem()

		obj.FieldByName("Id").SetInt(id)

		fmt.Println(t, obj, "+++++++++++", id, "++++++++", strings.Split(t1.String(), "."), "---------------")
		fmt.Printf("%T", t1)
	*/
	o := orm.NewOrm()

	m2m := o.QueryM2M(objrel, typeName)
	if m2m.Exist(obj) {
		b = true
		fmt.Println("obj Exist")
		return
	}

	return
}
