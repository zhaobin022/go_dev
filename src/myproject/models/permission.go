package models

import "github.com/astaxie/beego/orm"

type Permission struct {
	Id      int64
	Name    string `orm:"unique"`
	Comment string
	User    []*User `orm:"reverse(many)"`
	Role    []*Role `orm:"reverse(many)"`
}

func (r *Role) CheckRolePermission(permName string) (b bool) {
	b = false
	var permissions []*Permission
	o := orm.NewOrm()
	o.QueryTable(new(Permission)).Filter("Role__Role__Id", r.Id).All(&permissions)
	for _, permission := range permissions {
		if permName == permission.Name {
			b = true
			return
		}
	}
	return
}

func (u *User) IfhasPermisson(name string) (b bool) {
	b = false
	o := orm.NewOrm()

	var permission []*Permission

	num, err := o.QueryTable(new(Permission)).Filter("User__Id", u.Id).Filter("Name", name).All(&permission)
	if err != nil {
		return
	}

	if num > 0 {
		b = true
		return
	}

	var roles []*Role
	o.QueryTable(new(Role)).Filter("User__User__Id", u.Id).All(&roles)

	for _, g := range roles {
		b = g.CheckRolePermission(name)
		if b == true {
			return
		}

	}

	return
}
