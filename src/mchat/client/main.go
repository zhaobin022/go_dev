package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:10000")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	for {
		memu()
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		s = strings.Trim(s, "\r\n")
		switch s {
		case "1":
			err = register(conn)
			if err != nil {
				fmt.Println("login failed, err:", err)
				return
			}
		case "2":
			// login(conn)
			err = login(conn)
			if err != nil {
				fmt.Println("login failed, err:", err)
				return
			}
		default:
			fmt.Println("Please input the right option")
		}

	}

}
