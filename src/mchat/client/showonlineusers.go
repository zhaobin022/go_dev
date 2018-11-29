package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	. "mchat/model"
	. "mchat/util"
	"net"
)

func displayOnlineuser(conn net.Conn) {
	var msg Message
	msg.Cmd = DisplayOnlineUserCmd
	msg.Data = string(UserId)
	var buf [4]byte
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("format request online user data failed !")
		return
	}

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
		return
	}

	// msg, err = readPackage(conn)
	// if err != nil {
	// 	fmt.Println("read package failed, err:", err)
	// }

}

func displayUserList(data string) {
	fmt.Println(data)
	var user_map map[int]string = make(map[int]string)
	err := json.Unmarshal([]byte(data), &user_map)
	if err != nil {
		fmt.Println("user map json data unformat failed !", err)
		return
	}

	for k, v := range user_map {
		fmt.Println(k, " user : ", v, " login ")
	}

}
