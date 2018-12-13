package routers

import (
	"myproject/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var CheckLogin = func(ctx *context.Context) {
	loginUrl := beego.URLFor("LoginController.Get")

	login := ctx.Input.Session("login")
	if login != true {
		if ctx.Request.RequestURI != loginUrl {
			ctx.Redirect(302, loginUrl)
		}
	}
}

func init() {

	// beego.InsertFilter("/*", beego.BeforeRouter, CheckLogin)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/userlist", &controllers.UserApiController{})
	beego.Router("/useradd", &controllers.UserAddController{})
	beego.Router("/useredit/:id([0-9]+)", &controllers.UserEditController{})
	beego.Router("/changepass/:id([0-9]+)", &controllers.ChangePassController{})
}
