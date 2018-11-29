package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Student struct {
	Name string
	Sex  string
	Age  int
}

func (s *Student) Save() (err error) {

	data, err := json.Marshal(*s)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(`C:\Users\zhaobin\Desktop\go\src\example1\student.json`, data, 0777)
	return
}

func (s *Student) Load() (err error) {
	data, err := ioutil.ReadFile(`C:\Users\zhaobin\Desktop\go\src\example1\student.json`)
	if err != nil {
		fmt.Println("load file err : ", err)
	}
	err = json.Unmarshal(data, s)
	return

}

type Movie struct {
	Title string
	Year  int  `json:"released"`
	Color bool `json:"color, omitempty"`
}
