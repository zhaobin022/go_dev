package main

import "fmt"

func memu() {
	fmt.Println("1 . register ")
	fmt.Println("2 . login ")
}

func afterLoginMemu() {
	fmt.Println(`
	1 . list online user
	2 . talk to all
	`)
}
