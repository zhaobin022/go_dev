package controllers

import (
	"encoding/json"
	"fmt"
	. "myproject/models"
	"strconv"

	. "myproject/utils"

	"github.com/astaxie/beego/orm"
)

type MainController struct {
	BaseControl
}

func (this *MainController) Get() {
	this.Layout = "index.html"
	this.TplName = "index.html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["perm_menu"] = "perm_menu.html"
	this.LayoutSections["home"] = "home.html"
	this.LayoutSections["Scripts"] = "home_js.html"
}

type UserController struct {
	BaseControl
}

func (this *UserController) IsAjaxGet() {
	var pageSize int = 5
	pageStr := this.GetString("page")
	username := this.GetString("username")
	user := new(User)
	o := orm.NewOrm()
	qs := o.QueryTable(user)

	if username != "" {
		usernameCond := orm.NewCondition()

		qs = qs.SetCond(usernameCond.And("Name__icontains", username))
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	num, err := qs.Count()
	if err != nil {
		fmt.Println("get obj count error !", err)
	}
	var userPage *UserPage = &UserPage{}
	userPage.PaginatorMap = Paginator(page, pageSize, num)

	var currentPage, ok = (userPage.PaginatorMap)["currpage"].(int)
	if ok == false {
		fmt.Println("conver current page int failed !")
	}

	var userSlice []*User
	num, err = qs.OrderBy("Id").Limit(pageSize, (currentPage-1)*pageSize).All(&userSlice)
	if err != nil {
		fmt.Println(num, err)
	}
	userPage.UserSlice = userSlice
	this.Data["json"] = userPage
	this.ServeJSON()
}

func (this *UserController) Get() {
	isAjax := this.Ctx.Input.IsAjax()
	if isAjax {
		this.IsAjaxGet()
	} else {
		this.Data["IsUser"] = true
		this.Data["Perm"] = true
		this.Layout = "index.html"
		this.TplName = "index.html"
		this.LayoutSections = make(map[string]string)
		this.LayoutSections["home"] = "user.html"
		this.LayoutSections["Scripts"] = "user_scripts.html"
		this.LayoutSections["Css"] = "user_css.html"
	}

}

type UserPage struct {
	PaginatorMap map[string]interface{}
	UserSlice    []*User
}

func (this *UserController) Delete() {
	defer this.ServeJSON()
	data := this.Ctx.Input.RequestBody
	var ids []int
	err := json.Unmarshal(data, &ids)
	if err != nil {
		this.Data["json"] = err.Error()
		return
	}
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		fmt.Println(err)
	}

	var user = new(User)

	var relNameSlice []string = []string{"Permission", "Role"}
	err = DelObjAndRel(user, relNameSlice, &ids)

	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}

}

type UserAddController struct {
	BaseControl
}

func (this *UserAddController) Get() {
	o := orm.NewOrm()
	var permissions []*Permission
	permission := new(Permission)
	_, err := o.QueryTable(permission).All(&permissions)
	if err != nil {
		fmt.Println("get permission failed !")
	}

	var isAdmin map[string]int = make(map[string]int)
	isAdmin["是"] = 1
	isAdmin["否"] = 0

	var roles []*Role
	role := new(Role)
	_, err = o.QueryTable(role).All(&roles)
	if err != nil {
		fmt.Println("get role failed !")
	}

	this.Data["permissions"] = permissions
	this.Data["roles"] = roles
	this.Data["isAdmin"] = isAdmin
	this.TplName = "add_user.html"
}

type AddResponse struct {
	UserId int64
	Status bool
	Err    map[string]string
}

func (this *UserAddController) Post() {
	defer this.ServeJSON()
	res := &AddResponse{Err: make(map[string]string)}

	var userReq = new(UserAddRequest)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &userReq); err == nil {
		if userReq.Name == "" {
			res.Err["name"] = "用户名不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}

		if userReq.Password1 == "" || userReq.Password1 != userReq.Password2 {
			res.Err["all"] = "密码输入有问题!"
			res.Status = false
			this.Data["json"] = res
			return
		}

		user := new(User)
		o := orm.NewOrm()

		// 获取 QuerySeter 对象，user 为表名
		qs := o.QueryTable(user)
		exist := qs.Filter("Name", userReq.Name).Exist()

		if exist {
			res.Err["all"] = "用户已存在"
			res.Status = false
			this.Data["json"] = res
			return
		}

		user.Name = userReq.Name
		user.IsAdmin = userReq.IsAdmin
		// user.Permission = userReq.Permission

		encryPassFormt := GetEncryPass(userReq.Password1)
		user.Password = encryPassFormt
		userid, err := o.Insert(user)
		if err != nil {
			fmt.Println(err)
		}

		// 添加用户权限关系
		AddObjRel(user, userReq.Permission)
		// 添加用户角色
		AddObjRel(user, userReq.Role)

		if err != nil {
			res.Status = false
			res.Err["all"] = "添加用户失败"
			this.Data["json"] = res
			return
		} else {
			res.Status = true
			res.UserId = userid
			this.Data["json"] = res
			return
		}

	} else {
		res.Status = false
		res.Err["all"] = "转参错误"
		this.Data["json"] = res
		return
	}

	// this.TplName = "add_user.html"
}

type UserEditController struct {
	BaseControl
}

func (this *UserEditController) Get() {

	this.TplName = "edit_user.html"
	userIdStr := this.Ctx.Input.Param(":id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		fmt.Println("parse user id failed !")
		return
	}

	o := orm.NewOrm()
	user := new(User)
	user.Id = userId
	err = o.Read(user)
	if err != nil {
		fmt.Println("read user failed !", err)
		return
	}

	_, err = o.LoadRelated(user, "Role")
	if err != nil {
		fmt.Println("load user roles rel failed !")
	}

	_, err = o.LoadRelated(user, "Permission")
	if err != nil {
		fmt.Println("load user permission rel failed !")
	}

	var permissions []*Permission
	permission := new(Permission)
	o.QueryTable(permission).All(&permissions) // 返回 QuerySeter

	var roles []*Role
	role := new(Role)

	o.QueryTable(role).All(&roles) // 返回 QuerySeter

	var isAdmin map[string]bool = make(map[string]bool)
	isAdmin["是"] = true
	isAdmin["否"] = false
	this.Data["isAdmin"] = isAdmin
	this.Data["user"] = user
	this.Data["permissions"] = permissions
	this.Data["roles"] = roles

}

func (this *UserEditController) Post() {
	defer this.ServeJSON()

	res := &AddResponse{Err: make(map[string]string)}

	userStr := this.Ctx.Input.Param(":id")
	userId, err := strconv.ParseInt(userStr, 10, 64)
	if err != nil {
		res.Err["all"] = "用户Id接受错误!"
		res.Status = false
		this.Data["json"] = res
		return
	}
	var userReq = new(UserAddRequest)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &userReq); err == nil {

		user := new(User)
		o := orm.NewOrm()
		user.Id = userId

		// 获取 QuerySeter 对象，user 为表名
		err := o.Read(user)

		if err != nil {
			res.Err["all"] = "用户不存在"
			res.Status = false
			this.Data["json"] = res
			return
		}

		user.IsAdmin = userReq.IsAdmin
		//更新用户角色关系
		SyncObjRel(user, userReq.Role, "Role")
		//更新用户权限关系
		SyncObjRel(user, userReq.Permission, "Permission")

		if num, err := o.Update(user); err == nil {
			fmt.Println(num)
		}

		if err != nil {
			res.Status = false
			res.Err["all"] = "更改用户失败"
			this.Data["json"] = res
			return
		} else {
			res.Status = true
			res.UserId = 0
			this.Data["json"] = res
			return
		}

	} else {
		res.Status = false
		res.Err["all"] = "转参错误"
		this.Data["json"] = res
		return
	}
}

type ChangePassController struct {
	BaseControl
}

func (this *ChangePassController) Get() {
	this.TplName = "change_pass.html"
	resp := &Response{}
	userIdStr := this.Ctx.Input.Param(":id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		resp.Status = false
		resp.Msg = "parse user id failed !"
		return
	}
	o := orm.NewOrm()
	user := new(User)
	user.Id = userId
	err = o.Read(user)
	if err != nil {
		resp.Status = false
		resp.Msg = "read uesr error!"
		return
	}
	fmt.Println(user, 111)
	this.Data["user"] = user
}

func (this *ChangePassController) Put() {
	defer this.ServeJSON()
	resp := &Response{}
	var userReq = new(UserAddRequest)
	userIdStr := this.Ctx.Input.Param(":id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		resp.Status = false
		resp.Msg = "传参错误"
		this.Data["json"] = resp
		fmt.Println(resp.Msg)
		return
	}
	o := orm.NewOrm()
	user := new(User)
	user.Id = userId
	err = o.Read(user)
	if err != nil {
		resp.Status = false
		resp.Msg = "传参错误"
		this.Data["json"] = resp
		fmt.Println(resp.Msg)
		return
	}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &userReq); err != nil {

		resp.Status = false
		resp.Msg = "专参错误"
		this.Data["json"] = resp
		fmt.Println(resp.Msg)
		return
	}

	if userReq.Password1 == "" || userReq.Password2 == "" {
		resp.Status = false
		resp.Msg = "密码不能为空"
		this.Data["json"] = resp
		fmt.Println(resp.Msg)
		return
	}

	if userReq.Password1 != userReq.Password2 {
		resp.Status = false
		resp.Msg = "两次输入密码不同"
		fmt.Println(resp.Msg)
		this.Data["json"] = resp
		return
	}
	encryPassFormt := GetEncryPass(userReq.Password1)
	user.Password = encryPassFormt
	_, err = o.Update(user)
	if err != nil {
		resp.Status = false
		resp.Msg = "更新密码失败"
		fmt.Println(resp.Msg)
		this.Data["json"] = resp
		return
	}

	resp.Status = true
	resp.Msg = "更新密码成功"
	this.Data["json"] = resp
	return
}
