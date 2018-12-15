package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id         int64
	Name       string        `orm:"unique"`
	User       []*User       `orm:"rel(m2m)"`
	Permission []*Permission `orm:"rel(m2m)"`
}

func DelObjRel(obj interface{}, relObjSlice, reqRelObjSlice []interface{}) (err error) {
	// obj role relobj User Slice
	relObj := relObjSlice[0]
	objType := reflect.TypeOf(relObj)
	typeSlice := strings.Split(objType.String(), ".")
	typeName := typeSlice[len(typeSlice)-1]
	o := orm.NewOrm()
	if len(reqRelObjSlice) == 0 {
		m2m := o.QueryM2M(obj, typeName)
		nums, err := m2m.Clear()
		if err == nil {
			fmt.Println("Removed Tag Nums: ", nums)
		}
	} else {
		for _, tempobj := range relObjSlice {
			flag := false
			objValue := reflect.ValueOf(tempobj).Elem()
			objId := objValue.FieldByName("Id").Int()
			for _, reqObj := range reqRelObjSlice {
				reqObjValue := reflect.ValueOf(reqObj).Elem()
				reqObjId := reqObjValue.FieldByName("Id").Int()
				if objId == reqObjId {
					flag = true
					break
				}
			}
			if flag == false {
				m2m := o.QueryM2M(obj, typeName)
				num, err := m2m.Remove(tempobj)
				if err == nil {
					fmt.Println("Removed nums: ", num)
				}
			}
		}

	}

	/*
		_, err = o.LoadRelated(role, "User")
			if err != nil {
				fmt.Println("load rel user failed !")
			}
			for _, u := range role.User {
				flag := false
				for _, user := range roleReq.User {
					if user.Id == u.Id {
						flag = true
						break
					}
				}
				if flag == false {
					m2m := o.QueryM2M(role, "User")
					num, err := m2m.Remove(u)
					if err == nil {
						fmt.Println("Removed nums: ", num)
					}
				}
			}
	*/
	return
}

func AddObjRel(obj interface{}, relObjSlice []interface{}) (err error) {
	relObj := relObjSlice[0]
	objType := reflect.TypeOf(relObj)
	typeSlice := strings.Split(objType.String(), ".")
	typeName := typeSlice[len(typeSlice)-1]
	o := orm.NewOrm()
	m2m := o.QueryM2M(obj, typeName)
	for _, temprelobj := range relObjSlice {
		getValue := reflect.ValueOf(temprelobj).Elem()
		id := getValue.FieldByName("Id").Int()
		//create struct
		objStruct, err := GetObjFromStr(typeName)
		t := reflect.ValueOf(objStruct).Type()
		relobjPtr := reflect.New(t)
		relobj := relobjPtr.Elem()
		relobj.FieldByName("Id").SetInt(id)

		err = o.Read(relobjPtr.Interface())
		if err != nil {
			continue
		}

		if !m2m.Exist(relobjPtr.Interface()) {
			_, err = m2m.Add(relobjPtr.Interface())
			if err != nil {
				continue
			}
		}
	}
	return
}
