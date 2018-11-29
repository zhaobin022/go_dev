package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int, 2)
	// ch2 := make(chan int ,1)
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println("before put", i)
			ch1 <- i
			fmt.Println("after put", i)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			time.Sleep(time.Second * 5)
			s := <-ch1
			fmt.Println(s)
		}
	}()

	time.Sleep(time.Second * 100)
}
