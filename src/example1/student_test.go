package main

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	s := &Student{Name: "aaaa", Sex: "woman", Age: 111111}
	err := s.Save()

	if err != nil {
		t.Error("create student failed")
	}
	fmt.Println()
	return
}

func TestLoad(t *testing.T) {
	stu1 := &Student{Name: "aaaa", Sex: "woman", Age: 111111}
	stu2 := &Student{}
	err := stu2.Load()
	if err != nil {
		t.Error("load struct failed !")
	}

	if stu1.Name != stu2.Name {
		t.Error("name not  same")
	}

}
