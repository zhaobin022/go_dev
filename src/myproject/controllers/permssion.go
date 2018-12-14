package controllers

import (
	"encoding/json"
	"fmt"
	. "myproject/models"
	. "myproject/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type PermssionPage struct {
	PaginatorMap   map[string]interface{}
	PermssionSlice []*Permission
}

type PermissionController struct {
	BaseControl
}

func (this *PermissionController) IsAjaxGet() {
	defer this.ServeJSON()
	var pageSize int = 5
	pageStr := this.GetString("page")
	permName := this.GetString("permName")
	perm := new(Permission)
	o := orm.NewOrm()
	qs := o.QueryTable(perm)

	if permName != "" {
		permNameCond := orm.NewCondition()

		qs = qs.SetCond(permNameCond.And("Name__icontains", permName))
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
	permissionPage.PermssionSlice = PermssionSlice
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
	var permissionSlice []Permission
	_, err = o.QueryTable(permission).Filter("Id__in", ids).All(&permissionSlice)
	for _, p := range permissionSlice {
		m2mP := o.QueryM2M(&p, "User")
		nums, err := m2mP.Clear()
		if err == nil {
			fmt.Println("Removed permission Nums: ", nums)
		} else {
			fmt.Println(err)
			break
		}

		m2mR := o.QueryM2M(&p, "Role")
		nums, err = m2mR.Clear()
		if err == nil {
			fmt.Println("Removed role Nums: ", nums)
		} else {
			fmt.Println(err)
			break
		}

		o.Delete(&p)
	}

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

		// 添加用户权限关系
		m2m := o.QueryM2M(permission, "User")
		for _, user := range permissionReq.User {
			var u = &User{}
			u.Id = user.Id
			err := o.Read(u)
			if err != nil {
				continue
			}
			num, err := m2m.Add(u)
			if err != nil {
				fmt.Println("add user", num)
			}
		}
		// 添加用户角色

		for _, role := range permissionReq.Role {
			var r = &Role{}
			r.Id = role.Id
			err := o.Read(r)
			if err != nil {
				continue
			}
			m2mR := o.QueryM2M(permission, "Role")
			num, err := m2mR.Add(r)
			if err != nil {
				fmt.Println("add role", num)
			}
		}

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
	fmt.Println(permission, 11111111111111, this.TplName)

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
		if permissionReq.Name == "" {
			res.Msg = "权限名不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}

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

		permission.Name = permissionReq.Name
		for _, u := range permissionReq.User {
			err := o.Read(u)
			if err != nil {
				continue
			}

			m2m := o.QueryM2M(permission, "User")
			if !m2m.Exist(u) {
				m2m.Add(u)
			}
		}

		for _, r := range permissionReq.Role {
			err := o.Read(r)
			if err != nil {
				continue
			}

			m2m := o.QueryM2M(permission, "Role")
			if !m2m.Exist(r) {
				m2m.Add(r)
			}
		}
		_, err = o.LoadRelated(permission, "User")
		if err != nil {
			fmt.Println("load rel user failed !")
		}
		for _, u := range permission.User {
			flag := false
			for _, user := range permissionReq.User {
				if user.Id == u.Id {
					flag = true
					break
				}
			}
			if flag == false {
				m2m := o.QueryM2M(permission, "User")
				num, err := m2m.Remove(u)
				if err == nil {
					fmt.Println("Removed nums: ", num)
				}
			}
		}

		_, err = o.LoadRelated(permission, "Role")
		if err != nil {
			fmt.Println("load rel permission failed !")
		}

		for _, r := range permission.Role {
			flag := false
			for _, role := range permissionReq.Role {
				if r.Id == role.Id {
					flag = true
					break
				}
			}
			if flag == false {
				m2m := o.QueryM2M(permission, "Role")
				num, err := m2m.Remove(r)
				if err == nil {
					fmt.Println("Removed nums: ", num)
				}
			}
		}
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
