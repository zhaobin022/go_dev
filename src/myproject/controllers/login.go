package controllers

import (
	. "myproject/models"

	"github.com/astaxie/beego"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	username := this.GetString("username")
	password := this.GetString("password")
	user, err := LoginCheck(username, password)
	if err != nil {
		this.Data["Error"] = err
	} else {
		this.SetSession(LOGINSESSIONSTR, true)
		this.SetSession("username", user.Name)
		this.SetSession("userid", user.Id)
		this.Ctx.Redirect(302, "/")
	}
	this.TplName = "login.html"
}

type LogoutController struct {
	beego.Controller
}

func (this *LogoutController) Get() {
	this.DelSession(LOGINSESSIONSTR)
	loginUrl := beego.URLFor("LoginController.Get")

	this.Ctx.Redirect(302, loginUrl)
}
