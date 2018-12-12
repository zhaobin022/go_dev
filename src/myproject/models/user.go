package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/orm"
)

type UserAddRequest struct {
	User
	Password1 string
	Password2 string
}

type UserEditRequest struct {
	UserAddRequest
	page int
}

type User struct {
	Id         int64
	Name       string `orm:"unique"`
	Password   string `json:"-"`
	IsAdmin    bool
	Role       []*Role       `orm:"reverse(many)"` // 设置一对多的反向关系
	Permission []*Permission `orm:"rel(m2m)"`
}

func LoginCheck(username, password string) (user *User, e error) {

	if username == "" {
		e = fmt.Errorf("必须输入用户名")
		return
	}

	if password == "" {
		e = fmt.Errorf("必须输入密码")
		return
	}
	user = new(User)
	user.Name = username

	o := orm.NewOrm()
	err := o.QueryTable(new(User)).Filter("Name", username).One(user)
	if err != nil {
		e = fmt.Errorf("无此用户")
		return
	}

	encryPasswordSlice := strings.Split(user.Password, "|")
	salt, dbshasum := encryPasswordSlice[0], encryPasswordSlice[1]
	saltPassFormat := fmt.Sprintf("%s|%s", salt, password)
	passshasum := GenCryptPassStr(saltPassFormat)

	if passshasum == dbshasum {
		return
	} else {
		e = fmt.Errorf("密码不对!")
	}
	return
}
