package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.Debug = true

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterModel(new(User), new(Post), new(Profile), new(Tag))
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(10.12.9.195:3306)/test?charset=utf8")

	orm.RunSyncdb("default", false, true)
	fmt.Println("before register database !")
	// orm.RunSyncdb("default", false, true)
	fmt.Println("after register database !")

	// o := orm.NewOrm()
	// o.Using("default") // 默认使用 default，你可以指定为其他数据库

	// profile := new(Profile)
	// profile.Age = 30

	// user := new(User)
	// user.Profile = profile
	// user.Name = "slene"

	// fmt.Println(o.Insert(profile))
	// fmt.Println(o.Insert(user))

}

type User struct {
	Id      int64
	Name    string
	Profile *Profile `orm:"rel(one)"`      // OneToOne relation
	Post    []*Post  `orm:"reverse(many)"` // 设置一对多的反向关系
}

type Profile struct {
	Id   int
	Age  int16
	User *User `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}

type Post struct {
	Id    int
	Title string
	User  *User  `orm:"rel(fk)"` //设置一对多关系
	Tags  []*Tag `orm:"rel(m2m)"`
}

type Tag struct {
	Id    int
	Name  string
	Posts []*Post `orm:"reverse(many)"`
}
