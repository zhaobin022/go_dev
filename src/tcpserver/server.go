package main

import (
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("conn error : ", r)
		}
	}()

	for {
		const a = iota
		buf := make([]byte, 2048)

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read client content failed : ", err)
			return
		}

		content := string(buf[:n])
		fmt.Printf("msg from client %s : %s", conn.RemoteAddr, content)

	}
}

func main() {

	fmt.Println("staring server ")
	listen, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		fmt.Println("create listen failed : ", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("client connect error : ", err)
		}

		go process(conn)
	}
}
