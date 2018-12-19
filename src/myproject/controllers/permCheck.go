package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/astaxie/beego/context"

	. "myproject/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type PermDenyController struct {
	BaseControl
}

func (this *PermDenyController) Get() {
	this.Layout = "index.html"
	this.TplName = "index.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["perm_menu"] = "perm_menu.html"
	this.LayoutSections["home"] = "403.html"
	this.LayoutSections["Scripts"] = "403_js.html"
}

var CheckLogin = func(ctx *context.Context) {
	loginUrl := beego.URLFor("LoginController.Get")

	login := ctx.Input.Session("login")
	if login != true {
		fmt.Println(ctx.Request.RequestURI, loginUrl)
		if ctx.Request.RequestURI != loginUrl {
			ctx.Redirect(302, loginUrl)
		}
	}
}

func getUser(userid interface{}) (user *User, err error) {
	if userid == nil {
		err = errors.New("用户id为空")
		return
	}
	userId := userid.(int64)
	o := orm.NewOrm()
	user = new(User)
	user.Id = userId
	err = o.Read(user)
	if err != nil {
		return
	}
	return
}

func checkUrl(permission *Permission, uri, method string, tag bool) (b bool) {
	b = false
	urlSlice := strings.Split(permission.Url, ",")
	for _, url := range urlSlice {

		if tag == true {
			if url == uri {
				b = true
				return
			}
		} else if tag == false {
			if permission.Type == 1 {
				// pcreUrlSlice := strings.Split(permission.Url, ".")
				// m := strings.ToUpper(pcreUrlSlice[1])
				permUrl := beego.URLFor(strings.TrimSpace(url))
				fmt.Println("func : ", permUrl, uri)
				if permUrl == uri {
					b = true
				}
			} else if permission.Type == 2 {
				fmt.Println("url :", permission.Url, uri)
				if strings.TrimSpace(permission.Url) == uri {
					b = true
				}
			} else if permission.Type == 3 {
				match, _ := regexp.MatchString(permission.Url, uri)

				pcreUrlSlice := strings.Split(permission.Url, "|")
				fmt.Println("pcre match :", match)
				m := strings.ToUpper(pcreUrlSlice[1])
				if match && m == method {
					b = true
				}
			}
		}

		if b == true {
			return
		}
	}
	return
}

func DoPermCheck(user *User, uri, method string, tag bool) (b bool) {
	b = false
	if user.IsAdmin {
		b = true
		return
	}
	o := orm.NewOrm()

	_, err := o.LoadRelated(user, "Permission")
	if err != nil {
		return
	}
	fmt.Println(user.Permission, "permission")
	for _, perm := range user.Permission {
		b = checkUrl(perm, uri, method, tag)
		if b == true {
			return
		}
	}

	_, err = o.LoadRelated(user, "Role")
	if err != nil {
		return
	}
	for _, r := range user.Role {
		_, err = o.LoadRelated(r, "Permission")
		if err != nil {
			return
		}
		for _, p := range r.Permission {
			b = checkUrl(p, uri, method, tag)
			if b == true {
				return
			}
		}
	}

	return
}
