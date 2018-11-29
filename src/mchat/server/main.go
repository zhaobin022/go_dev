package main

import (
	"fmt"
	_ "mchat/util"
)

var CliMgr ClientMgr = ClientMgr{}

func init() {

	CliMgr.allClientMap = make(map[int]Client)
}

func main() {
	fmt.Println("starting server")
	runServer("0.0.0.0:10000")
}
