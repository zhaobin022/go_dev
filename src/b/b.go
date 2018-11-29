package main

import (
	"fmt"
)

type Student struct {
	name string
}

func modify(s Student) {
	s.name = "modify"
	fmt.Printf("%p\n", &s)
}
func main() {

	var s1 Student = Student{name: "old"}
	fmt.Printf("%p\n", &s1)

	modify(s1)
	fmt.Printf("%p\n", &s1)

}
