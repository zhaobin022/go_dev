package controllers

import (
	"errors"
	"fmt"
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

func DoPermCheck(user *User, uri string) (b bool) {
	b = false
	if user.IsAdmin {
		b = true
		return
	}
	o := orm.NewOrm()

	var permissions []*Permission
	var permission = new(Permission)
	_, err := o.QueryTable(permission).Filter("User__User__Id", user.Id).All(&permissions)
	if err != nil {
		return
	}

	for _, perm := range permissions {

		ctrlSlice := strings.Split(perm.Url, ",")
		for _, controller := range ctrlSlice {
			permUrl := beego.URLFor(strings.TrimSpace(controller))
			fmt.Println(permUrl, uri, "ppppppppppppppppppppppppp")
			if permUrl == uri {
				b = true
				return
			}
		}
	}
	role := new(Role)
	var roles []*Role
	_, err = o.QueryTable(role).Filter("User__User__Id", user.Id).All(&roles)
	if err != nil {
		return
	}
	for _, r := range roles {

		_, err = o.QueryTable(permission).Filter("Role__Role__Id", r.Id).All(&permissions)
		for _, p := range permissions {
			ctrlSlice := strings.Split(p.Url, ",")
			for _, controller := range ctrlSlice {
				permUrl := beego.URLFor(strings.TrimSpace(controller))
				if permUrl == uri {
					b = true
					return
				}
			}
		}
	}

	return
}
