package main

import (
	"fmt"
)

type Instance struct {
	host string
	port int
}

func CreateInstance(host string, port int) Instance {
	return Instance{host: host, port: port}
}

func (i *Instance) String() {
	fmt.Printf("%s:%d", i.host, i.port)
}
