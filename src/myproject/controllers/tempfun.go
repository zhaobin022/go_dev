package controllers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

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

func PermCtl(ctl *BaseControl, permStr string) (b bool) {
	method := ctl.Ctx.Request.Method
	b = false
	ctx := ctl.Ctx
	userId := ctx.Input.Session("userid")
	// uri := beego.URLFor(permStr)
	user, err := getUser(userId)
	if err != nil {
		return
	}
	b = DoPermCheck(user, permStr, method, true)
	// fmt.Println(user, uri, b, permStr, "html permission check")
	return
}
