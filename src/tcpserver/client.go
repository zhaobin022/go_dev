package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:5000")

	if err != nil {
		fmt.Println("dial server error : ", err)
	}

	// var buf = make([]byte, 2048)
	inputReader := bufio.NewReader(os.Stdin)
	for {
		s, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println("read from input error : ", err)
		}
		// fmt.Printf("-%s-\n", s)
		if strings.Trim(s, "\r\n") == "q" {
			break
		}
		_, err = conn.Write([]byte(s))
		if err != nil {
			fmt.Println("send to server error :", err)
		}

	}

}
