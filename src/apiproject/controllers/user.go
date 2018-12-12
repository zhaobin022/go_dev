package controllers

import (
	"apiproject/models"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var (
	o orm.Ormer
)

func init() {

	o = orm.NewOrm()
	o.Using("default")
	orm.RunSyncdb("default", false, true)
}

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {

	var user *models.User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	fmt.Println(string(u.Ctx.Input.RequestBody))
	fmt.Println(user)

	fmt.Println("----------------------------------", err)
	fmt.Println(o.Insert(user))
	u.Data["json"] = map[string]string{"uid": string(user.Id)}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {

	// profile := new(models.Profile)
	// profile.Age = 30

	// user := new(models.User)
	// user.Profile = profile
	// user.Name = "slene"

	// fmt.Println(o.Insert(profile))
	// fmt.Println(o.Insert(user))
	var userList []models.User
	users := o.QueryTable(new(models.User))
	users.All(&userList)

	u.Data["json"] = userList
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid, err := u.GetInt64(":uid")
	fmt.Println(uid, "------------------------")
	// userId, err := strconv.Atoi(uid)
	// userId, err := strconv.ParseInt(uid, 10, 64)
	user := models.User{Id: uid}
	err = o.Read(&user)

	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
		u.Data["json"] = err.Error()
	} else if err == orm.ErrMissPK {
		fmt.Println("找不到主键")
		u.Data["json"] = err.Error()
	} else {
		fmt.Println(user.Id, user.Name, "=====================")
		u.Data["json"] = user
	}

	// if err != nil {
	// 	logs.Error("get user id error")
	// 	u.Data["json"] = err.Error()
	// 	return
	// }

	// user := &models.User{Id: userId}
	// if uid != "" {
	// 	user, err := o.Read(user, ["UserName"]string)
	// 	if err != nil {
	// 		u.Data["json"] = err.Error()
	// 	} else {
	// 		u.Data["json"] = user
	// 	}
	// }
	u.ServeJSON()
}

func (u *UserController) Delete() {
	uid, err := u.GetInt64(":uid")
	user := new(models.User)
	// userId, err := strconv.Atoi(uid)
	// userId, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		u.Data["json"] = err.Error()
	}
	user.Id = uid
	o.Delete(user)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}
