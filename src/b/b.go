package main

import (
	"fmt"
	"regexp"
	// "WebSocketInBeego/models"
)

// var engine *xorm.Engine

func main() {
	text := `Hello 世界！123 Go.`

	// 查找连续的小写字母
	reg := regexp.MustCompile(`(\w)e(\w)`)
	fmt.Printf("%q\n", reg.FindAllStringSubmatch(text, -1))
}
