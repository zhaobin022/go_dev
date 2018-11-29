package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	. "mchat/model"
	. "mchat/util"
	"net"
	"os"
	"strings"
)

func sendMsg(conn net.Conn, m string) (err error) {
	var msg Message
	msg.Cmd = TalkToallCmd
	msg.Data = m
	data, err := json.Marshal(msg)
	var buf [4]byte
	packLen := uint32(len(data))

	// fmt.Println("packlen:", packLen)
	binary.BigEndian.PutUint32(buf[0:4], packLen)

	n, err := conn.Write(buf[:])
	if err != nil || n != 4 {
		fmt.Println("write data  failed")
		return
	}

	_, err = conn.Write([]byte(data))
	if err != nil {
		fmt.Println("send msg failed")
		return
	}
	return
}
func talkToall(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input from termnial failed !")
			continue
		}
		s = strings.Trim(s, "\r\n")
		if len(s) == 0 {
			continue
		}

		err = sendMsg(conn, s)
	}
}

func displayTalkMsg(conn net.Conn, msg Message) {
	var data_list []string
	err := json.Unmarshal([]byte(msg.Data), &data_list)
	if err != nil {
		fmt.Println("unformat user data error ")
		return
	}
	for index, v := range data_list {
		data_list[index] = strings.Trim(v, "\r\n")
	}
	fmt.Printf("from %s : %s \r\n", data_list[0], data_list[1])

}
