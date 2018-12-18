package controllers

import (
	"fmt"
	. "myproject/models"
	"regexp"

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
	PermControl
}

type PermControl struct {
}

func (this *BaseControl) Prepare() {

	var username = this.GetSession("username")
	if username != "" {
		this.Data["username"] = username
	}

	this.CheckPerm()
}

func (this *BaseControl) CheckPerm() {
	ctx := this.Ctx
	userId := ctx.Input.Session("userid")
	permDenyUrl := beego.URLFor("PermDenyController.Get")
	loginUrl := beego.URLFor("LoginController.Get")
	user, err := getUser(userId)

	if err != nil {
		fmt.Println(ctx.Request.RequestURI, loginUrl, permDenyUrl)
		if ctx.Request.RequestURI != loginUrl && ctx.Request.RequestURI != permDenyUrl {
			if ctx.Input.IsAjax() {
				var basePage *BasePage = &BasePage{}
				basePage.PermDeny = true
				this.Data["json"] = basePage
				this.ServeJSON()
			} else {
				ctx.Redirect(302, permDenyUrl)
			}
		}
		return
	}

	var uri string
	exp3 := regexp.MustCompile(`(.*)\?.*`)
	fmt.Println(uri, "888888888")
	result3 := exp3.FindAllStringSubmatch(ctx.Input.URI(), -1)
	if len(result3) > 0 {
		fmt.Println(result3[0], "999999999999999")
		uri = result3[0][1]
	} else {
		uri = ctx.Request.RequestURI
	}

	fmt.Println(uri, "------------------", ctx.Input.URI())
	ok := DoPermCheck(user, uri)
	fmt.Println(ok, "___________________________")
	if ok == false {
		if ctx.Request.RequestURI != loginUrl && ctx.Request.RequestURI != permDenyUrl {
			if ctx.Input.IsAjax() {
				var basePage *BasePage = &BasePage{}
				basePage.PermDeny = true
				this.Data["json"] = basePage
				this.ServeJSON()
			} else {
				ctx.Redirect(302, permDenyUrl)
			}
		}
	}
	return
}

func init() {
	beego.AddFuncMap("IfObjInObjRel", IfObjInObjRel)
}
