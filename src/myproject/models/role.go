package models

type Role struct {
	Id         int64
	Name       string        `orm:"unique"`
	User       []*User       `orm:"rel(m2m)"`
	Permission []*Permission `orm:"rel(m2m)"`
}
