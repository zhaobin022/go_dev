package controllers

import (
	. "myproject/models"
)

type BasePage struct {
	PaginatorMap map[string]interface{}
}

type UserPage struct {
	BasePage
	ObjSlice []*User
}

type RolePage struct {
	BasePage
	ObjSlice []*Role
}

type PermssionPage struct {
	BasePage
	ObjSlice []*Permission
}
