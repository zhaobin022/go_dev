package controllers

import (
	"encoding/json"
	"fmt"
	. "myproject/models"
	. "myproject/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type PermissionController struct {
	BaseControl
}

func (this *PermissionController) IsAjaxGet() {
	defer this.ServeJSON()
	var pageSize int = 5
	pageStr := this.GetString("page")
	permName := this.GetString("Name")
	url := this.GetString("Url")

	perm := new(Permission)
	o := orm.NewOrm()
	qs := o.QueryTable(perm)

	if permName != "" {
		permNameCond := orm.NewCondition()
		qs = qs.SetCond(permNameCond.And("Name__icontains", permName))
	}

	if url != "" {
		urlCond := orm.NewCondition()
		qs = qs.SetCond(urlCond.And("Url__icontains", url))
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	num, err := qs.Count()
	if err != nil {
		fmt.Println("get obj count error !", err)
	}
	var permissionPage *PermssionPage = &PermssionPage{}
	permissionPage.PaginatorMap = Paginator(page, pageSize, num)

	var currentPage, ok = (permissionPage.PaginatorMap)["currpage"].(int)
	if ok == false {
		fmt.Println("conver current page int failed !")
	}

	var PermssionSlice []*Permission
	num, err = qs.OrderBy("Id").Limit(pageSize, (currentPage-1)*pageSize).All(&PermssionSlice)
	if err != nil {
		fmt.Println(num, err)
	}
	permissionPage.ObjSlice = PermssionSlice
	this.Data["json"] = permissionPage
}

func (this *PermissionController) Get() {
	isAjax := this.Ctx.Input.IsAjax()
	if isAjax {
		this.IsAjaxGet()
	} else {
		this.Data["IsPermssion"] = true
		this.Data["Perm"] = true
		this.Layout = "index.html"
		this.TplName = "index.html"
		this.LayoutSections = make(map[string]string)
		this.LayoutSections["perm_menu"] = "perm_menu.html"
		this.LayoutSections["home"] = "perm.html"
		this.LayoutSections["Scripts"] = "perm_scripts.html"
		this.LayoutSections["Css"] = "perm_css.html"
	}
}

func (this *PermissionController) Delete() {
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

	var permission = new(Permission)

	var relNameSlice []string = []string{"User", "Role"}
	err = DelObjAndRel(permission, relNameSlice, &ids)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}

}

type PermssionAddController struct {
	BaseControl
}

func (this *PermssionAddController) Get() {
	this.TplName = "add_perm.html"
	o := orm.NewOrm()
	var roles []*Role
	role := new(Role)
	_, err := o.QueryTable(role).All(&roles)
	if err != nil {
		fmt.Println("get permission failed !")
	}

	var users []*User
	user := new(User)
	_, err = o.QueryTable(user).All(&users)
	if err != nil {
		fmt.Println("get user failed !")
	}

	this.Data["roles"] = roles

	this.Data["users"] = users
}

func (this *PermssionAddController) Post() {
	defer this.ServeJSON()
	res := new(Response)

	var permissionReq = new(Permission)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &permissionReq); err == nil {
		if permissionReq.Name == "" {
			res.Msg = "权限名不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}

		permission := new(Permission)

		o := orm.NewOrm()

		// 获取 QuerySeter 对象，user 为表名
		qs := o.QueryTable(permission)
		exist := qs.Filter("Name", permissionReq.Name).Exist()
		if exist {
			res.Msg = "权限已存在"
			res.Status = false
			this.Data["json"] = res
			return
		}
		permission.Name = permissionReq.Name

		exist = qs.Filter("Url", permissionReq.Url).Exist()
		if exist {
			res.Msg = "URL已存在"
			res.Status = false
			this.Data["json"] = res
			return
		}

		permission.Url = permissionReq.Url

		_, err = o.Insert(permission)
		if err != nil {
			fmt.Println(err)
		}

		// 添加权限用户关系
		AddObjRel(permission, permissionReq.User)
		// 添加权限角色关系
		AddObjRel(permission, permissionReq.Role)

		if err != nil {
			res.Status = false
			res.Msg = "添加权限失败"
			this.Data["json"] = res
			return
		} else {
			res.Status = true
			res.Msg = "添加权限成功"
			this.Data["json"] = res
			return
		}

	} else {
		res.Status = false
		res.Msg = "传参错误"
		this.Data["json"] = res
		return
	}
}

type PermissionEditController struct {
	BaseControl
}

func (this *PermissionEditController) Get() {
	this.TplName = "edit_perm.html"
	permissionIdStr := this.Ctx.Input.Param(":id")
	permissionId, err := strconv.ParseInt(permissionIdStr, 10, 64)

	if err != nil {
		fmt.Println("parse role id failed !")
		return
	}

	o := orm.NewOrm()
	permission := new(Permission)
	permission.Id = permissionId
	err = o.Read(permission)
	if err != nil {
		fmt.Println("read role failed !", err)
		return
	}

	_, err = o.LoadRelated(permission, "User")
	if err != nil {
		fmt.Println("load user user rel failed !")
	}

	_, err = o.LoadRelated(permission, "Role")
	if err != nil {
		fmt.Println("load user orle rel failed !")
	}

	var roles []*Role
	role := new(Role)
	o.QueryTable(role).All(&roles) // 返回 QuerySeter

	var users []*User
	user := new(User)
	o.QueryTable(user).All(&users) // 返回 QuerySeter

	this.Data["permission"] = permission
	this.Data["roles"] = roles
	this.Data["users"] = users
	this.Data["urlType"] = PermUrlType
	
}

func (this *PermissionEditController) Post() {
	defer this.ServeJSON()

	res := &Response{}

	permStr := this.Ctx.Input.Param(":id")
	permssionId, err := strconv.ParseInt(permStr, 10, 64)
	if err != nil {
		res.Msg = "权限Id接受错误!"
		res.Status = false
		this.Data["json"] = res
		return
	}
	var permissionReq = new(Permission)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &permissionReq); err == nil {
		// if permissionReq.Name == "" {
		// 	res.Msg = "权限名不能为空!"
		// 	res.Status = false
		// 	this.Data["json"] = res
		// 	return
		// }

		if permissionReq.Url == "" {
			res.Msg = "权限名URL不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}
		permission := new(Permission)
		o := orm.NewOrm()
		permission.Id = permssionId

		// 获取 QuerySeter 对象，user 为表名
		err := o.Read(permission)
		if err != nil {
			res.Msg = "权限不存在"
			res.Status = false
			this.Data["json"] = res
			return
		}

		err = o.Begin()
		if err != nil {
			res.Status = false
			res.Msg = "更新失败"
			this.Data["json"] = res
			return
		}

		permission.Url = permissionReq.Url

		//更新用户角色关系
		err = SyncObjRel(permission, permissionReq.Role, "Role")
		//更新用户权限关系
		err = SyncObjRel(permission, permissionReq.User, "User")

		if num, err := o.Update(permission); err == nil {
			fmt.Println(num)
		}

		if err != nil {
			err = o.Rollback()
			res.Status = false
			res.Msg = "更改权限失败"
			this.Data["json"] = res
			return
		} else {
			err = o.Commit()
			res.Status = true
			res.Msg = "更新权限成功"
			this.Data["json"] = res
			return
		}

	} else {
		res.Status = false
		res.Msg = "转参错误"
		this.Data["json"] = res
		return
	}

}
