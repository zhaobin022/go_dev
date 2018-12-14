package main

import (
	_ "myproject/routers"

	_ "myproject/code"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
