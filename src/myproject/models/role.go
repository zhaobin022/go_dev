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

/*
func getStructName(obj interface{}) {
	objType := reflect.TypeOf(obj)
	typeSlice := strings.Split(objType.String(), ".")
	typeName := typeSlice[len(typeSlice)-1]
}
*/

func ClearObjRel(obj interface{}, structName string) {
	o := orm.NewOrm()
	m2m := o.QueryM2M(obj, structName)
	nums, err := m2m.Clear()
	if err == nil {
		fmt.Println("Removed obj rel: ", nums)
	}
}

func SyncObjRel(obj, reqRelObj interface{}, structName string) (err error) {

	o := orm.NewOrm()
	_, err = o.LoadRelated(obj, structName)
	if err != nil {
		fmt.Println("load rel obj failed !")
	}
	objValue := reflect.ValueOf(obj).Elem()
	relObjSliceValue := objValue.FieldByName(structName)
	reqRelObjSliceValue := reflect.ValueOf(reqRelObj)
	if relObjSliceValue.Len() == 0 && reqRelObjSliceValue.Len() == 0 {
		return
	}

	if relObjSliceValue.Len() > 0 && reqRelObjSliceValue.Len() == 0 {
		ClearObjRel(obj, structName)
		return
	}

	AddObjRel(obj, reqRelObj)

	// m2m := o.QueryM2M(obj, typeName)

	relObjSlice := relObjSliceValue.Slice(0, relObjSliceValue.Len())
	reqRelObjSlice := reqRelObjSliceValue.Slice(0, reqRelObjSliceValue.Len())

	for i := 0; i < relObjSlice.Len(); i++ {
		var flag = false
		v := relObjSlice.Index(i)
		id := v.Elem().FieldByName("Id").Int()
		for j := 0; j < reqRelObjSlice.Len(); j++ {
			z := reqRelObjSlice.Index(j)
			reqId := z.Elem().FieldByName("Id").Int()
			if id == reqId {
				flag = true
				break
			}
		}

		if flag == false {
			m2m := o.QueryM2M(obj, structName)
			num, err := m2m.Remove(v.Interface())
			if err == nil {
				fmt.Println("Removed nums: ", num)
			}
		}

	}
	return
}

func AddObjRel(obj, relObjSlice interface{}) (err error) {
	o := orm.NewOrm()

	getValue := reflect.ValueOf(relObjSlice)
	if getValue.Len() == 0 {
		return
	}

	slice := getValue.Slice(0, getValue.Len())

	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		//获取对象名
		objType := v.Type()
		typeSlice := strings.Split(objType.String(), ".")
		typeName := typeSlice[len(typeSlice)-1]

		id := v.Elem().FieldByName("Id").Int()
		//create struct
		objStruct, err := GetObjFromStr(typeName)
		t := reflect.ValueOf(objStruct).Type()
		relobjPtr := reflect.New(t)
		relobj := relobjPtr.Elem()
		relobj.FieldByName("Id").SetInt(id)
		//效验实例
		err = o.Read(relobjPtr.Interface())
		if err != nil {
			continue
		}
		//更新关系表
		m2m := o.QueryM2M(obj, typeName)
		if !m2m.Exist(relobjPtr.Interface()) {
			_, err = m2m.Add(relobjPtr.Interface())
			if err != nil {
				continue
			}
		}
	}
	return
}

func DelObjAndRel(obj interface{}, relNameSlice []string, ids *[]int) (err error) {

	objType := reflect.TypeOf(obj)
	objSlice := reflect.MakeSlice(reflect.SliceOf(objType), 0, 0)
	x := reflect.New(objSlice.Type())
	x.Elem().Set(objSlice)
	o := orm.NewOrm()
	_, err = o.QueryTable(obj).Filter("Id__in", ids).All(x.Interface())
	if err != nil {
		return
	}

	slice := x.Elem().Slice(0, x.Elem().Len())

	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)
		for _, name := range relNameSlice {
			m2m := o.QueryM2M(v.Interface(), name)
			nums, err := m2m.Clear()
			if err == nil {
				fmt.Println("Removed obj Nums: ", nums)
			} else {
				fmt.Println(err)
				return err
			}

		}
		_, err = o.Delete(v.Interface())
		if err != nil {
			return
		}
	}
	return
}
