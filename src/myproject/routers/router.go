package routers

import (
	"myproject/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.InsertFilter("/*", beego.BeforeRouter, controllers.CheckLogin)
	// beego.InsertFilter("/*", beego.BeforeRouter, controllers.CheckPerm)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/logout", &controllers.LogoutController{})

	beego.Router("/403", &controllers.PermDenyController{})

	// 用户管理路由
	beego.Router("/user", &controllers.UserController{})
	beego.Router("/useradd", &controllers.UserAddController{})
	beego.Router("/useredit/:id([0-9]+)", &controllers.UserEditController{})
	beego.Router("/changepass/:id([0-9]+)", &controllers.ChangePassController{})
	//角色管理路由
	beego.Router("/role", &controllers.RoleController{})
	beego.Router("/roleadd", &controllers.RoleAddController{})
	beego.Router("/roleedit/:id([0-9]+)", &controllers.RoleEditController{})
	//权限管理路由
	beego.Router("/permission", &controllers.PermissionController{})
	beego.Router("/permissionadd", &controllers.PermssionAddController{})
	beego.Router("/permissionedit/:id([0-9]+)", &controllers.PermissionEditController{})

}
