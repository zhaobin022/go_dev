package controllers

import (
	"encoding/json"
	"fmt"
	. "myproject/models"
	. "myproject/utils"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type RoleController struct {
	BaseControl
}

func (this *RoleController) IsAjaxGet() {
	defer this.ServeJSON()
	var pageSize int = 5
	pageStr := this.GetString("page")
	rolename := this.GetString("rolename")
	role := new(Role)
	o := orm.NewOrm()
	qs := o.QueryTable(role)

	if rolename != "" {
		rolenameCond := orm.NewCondition()

		qs = qs.SetCond(rolenameCond.And("Name__icontains", rolename))
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	num, err := qs.Count()
	if err != nil {
		fmt.Println("get obj count error !", err)
	}
	var rolePage *RolePage = &RolePage{}
	rolePage.PaginatorMap = Paginator(page, pageSize, num)

	var currentPage, ok = (rolePage.PaginatorMap)["currpage"].(int)
	if ok == false {
		fmt.Println("conver current page int failed !")
	}

	var roleSlice []*Role
	num, err = qs.OrderBy("Id").Limit(pageSize, (currentPage-1)*pageSize).All(&roleSlice)
	if err != nil {
		fmt.Println(num, err)
	}
	rolePage.RoleSlice = roleSlice
	this.Data["json"] = rolePage

}

func (this *RoleController) Get() {
	isAjax := this.Ctx.Input.IsAjax()
	if isAjax {
		this.IsAjaxGet()
	} else {
		this.Data["IsRole"] = true
		this.Data["Perm"] = true
		this.Layout = "index.html"
		this.TplName = "index.html"
		this.LayoutSections = make(map[string]string)
		// this.LayoutSections["perm_menu"] = "perm_menu.html"
		this.LayoutSections["home"] = "role.html"
		this.LayoutSections["Scripts"] = "role_scripts.html"
		this.LayoutSections["Css"] = "role_css.html"
	}
}

func (this *RoleController) Delete() {
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

	var role = new(Role)
	// var roleSlice []Role
	var relNameSlice []string = []string{"Permission", "User"}
	err = DelObjAndRel(role, &relNameSlice, &ids)
	/*
		_, err = o.QueryTable(role).Filter("Id__in", ids).All(&roleSlice)
		for _, r := range roleSlice {
			m2mP := o.QueryM2M(&r, "Permission")
			nums, err := m2mP.Clear()
			if err == nil {
				fmt.Println("Removed permission Nums: ", nums)
			} else {
				fmt.Println(err)
				break
			}

			m2mR := o.QueryM2M(&r, "User")
			nums, err = m2mR.Clear()
			if err == nil {
				fmt.Println("Removed user Nums: ", nums)
			} else {
				fmt.Println(err)
				break
			}

			o.Delete(&r)
		}
	*/
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
}

type RolePage struct {
	PaginatorMap map[string]interface{}
	RoleSlice    []*Role
}

type RoleAddController struct {
	BaseControl
}

func (this *RoleAddController) Get() {
	this.TplName = "add_role.html"
	o := orm.NewOrm()
	var permissions []*Permission
	permission := new(Permission)
	_, err := o.QueryTable(permission).All(&permissions)
	if err != nil {
		fmt.Println("get permission failed !")
	}

	var users []*User
	user := new(User)
	_, err = o.QueryTable(user).All(&users)
	if err != nil {
		fmt.Println("get user failed !")
	}

	this.Data["permissions"] = permissions

	this.Data["users"] = users
}

func (this *RoleAddController) Post() {
	defer this.ServeJSON()
	res := new(Response)

	var roleReq = new(Role)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &roleReq); err == nil {
		if roleReq.Name == "" {
			res.Msg = "角色名不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}

		role := new(Role)

		o := orm.NewOrm()

		// 获取 QuerySeter 对象，user 为表名
		qs := o.QueryTable(role)
		count, err := qs.Filter("Name", roleReq.Name).Count()
		if err != nil {
			fmt.Println("query count error ", err)
		}

		if count > 0 {
			res.Msg = "角色已存在"
			res.Status = false
			this.Data["json"] = res
			return
		}
		role.Name = roleReq.Name
		_, err = o.Insert(role)
		if err != nil {
			fmt.Println(err)
		}

		// 添加用户权限关系
		AddObjRel(role, roleReq.Permission)
		// 添加用户角色
		AddObjRel(role, roleReq.User)

		if err != nil {
			res.Status = false
			res.Msg = "添加角色失败"
			this.Data["json"] = res
			return
		} else {
			res.Status = true
			res.Msg = "添加角色成功"
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

type RoleEditController struct {
	BaseControl
}

func (this *RoleEditController) Get() {
	this.TplName = "edit_role.html"
	roleIdStr := this.Ctx.Input.Param(":id")
	roleId, err := strconv.ParseInt(roleIdStr, 10, 64)

	if err != nil {
		fmt.Println("parse role id failed !")
		return
	}

	o := orm.NewOrm()
	role := new(Role)
	role.Id = roleId
	err = o.Read(role)
	if err != nil {
		fmt.Println("read role failed !", err)
		return
	}

	_, err = o.LoadRelated(role, "User")
	if err != nil {
		fmt.Println("load user roles rel failed !")
	}

	_, err = o.LoadRelated(role, "Permission")
	if err != nil {
		fmt.Println("load user permission rel failed !")
	}

	var permissions []*Permission
	permission := new(Permission)
	o.QueryTable(permission).All(&permissions) // 返回 QuerySeter

	var users []*User
	user := new(User)
	o.QueryTable(user).All(&users) // 返回 QuerySeter

	this.Data["role"] = role
	this.Data["permissions"] = permissions
	this.Data["users"] = users

}

func (this *RoleEditController) Post() {
	defer this.ServeJSON()

	res := &Response{}

	roleStr := this.Ctx.Input.Param(":id")
	roleId, err := strconv.ParseInt(roleStr, 10, 64)
	if err != nil {
		res.Msg = "角色Id接受错误!"
		res.Status = false
		this.Data["json"] = res
		return
	}
	var roleReq = new(Role)
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &roleReq); err == nil {
		if roleReq.Name == "" {
			res.Msg = "角色名不能为空!"
			res.Status = false
			this.Data["json"] = res
			return
		}
		role := new(Role)
		o := orm.NewOrm()
		role.Id = roleId

		// 获取 QuerySeter 对象，user 为表名
		err := o.Read(role)
		if err != nil {
			res.Msg = "角色不存在"
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

		role.Name = roleReq.Name

		SyncObjRel(role, roleReq.User, "User")
		SyncObjRel(role, roleReq.Permission, "Permission")

		if num, err := o.Update(role); err == nil {
			fmt.Println(num)
		}

		if err != nil {
			err = o.Rollback()
			res.Status = false
			res.Msg = "更改用户失败"
			this.Data["json"] = res
			return
		} else {
			err = o.Commit()
			res.Status = true
			res.Msg = "更新用户成功"
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
