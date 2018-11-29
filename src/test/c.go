package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {

	var ins []Instance
	for i := 0; i < 4; i++ {
		in := CreateInstance(fmt.Sprintf("192.168.11.%d", i), 8000+rand.Intn(1000))
		ins = append(ins, in)
	}
	var balanceName string
	if len(os.Args) != 2 {
		fmt.Println("params not enough")
		return
	}
	balanceName = os.Args[1]

	bai, err := GetBalance(balanceName)
	if err != nil {
		fmt.Println("error : ", err)
		return
	}
	var b Balance
	b = bai
	for {
		in, err := b.DoBalance(ins)
		if err != nil {
			break
		}
		fmt.Println(in)
		time.Sleep(time.Second)
	}

}
