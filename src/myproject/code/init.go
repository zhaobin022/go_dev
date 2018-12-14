package code

import (
	"fmt"
)

const (
	UserNameEmpty       = iota + 1000 // 0
	UserPasswordEmpty                 // 1
	UserPasswordNotSame               // 2
	UserPasswordTooShot               // 3
)

type Code struct {
	id  int
	msg string
}

var UserCodeSet map[int]*Code

func initError() {
	UserCodeSet = make(map[int]*Code)
	UserCodeSet[UserNameEmpty] = &Code{id: UserNameEmpty, msg: "用户名不能为空"}
}

func init() {
	initError()
	fmt.Println(UserCodeSet)
}
