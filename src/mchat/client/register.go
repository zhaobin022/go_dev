package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"mchat/model"
	"mchat/util"
	"net"
	"os"
	"strconv"
	"strings"
)

func register(conn net.Conn) (err error) {
	var msg model.Message
	msg.Cmd = util.UserRegister

	var registerCmd model.RegisterCmd
	for {
		fmt.Println("Plase input the user id")
		reader := bufio.NewReader(os.Stdin)
		userid, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input userid error")
			continue
		}

		userid = strings.Trim(userid, "\r\n")
		id, err := strconv.Atoi(userid)
		if err != nil {
			fmt.Println("must input number !")
			continue
		}

		fmt.Println("Plase input password")
		// reader = bufio.NewReader(os.Stdin)
		passwd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input password error")
			continue
		}

		passwd = strings.Trim(passwd, "\r\n")

		fmt.Println("Plase input the nick name")
		// reader = bufio.NewReader(os.Stdin)
		nick, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input nick name error")
			continue
		}

		registerCmd.User.UserId = id
		registerCmd.User.Passwd = passwd
		registerCmd.User.Nick = nick
		break
	}

	data, err := json.Marshal(registerCmd)
	if err != nil {
		return
	}

	msg.Data = string(data)
	data, err = json.Marshal(msg)
	if err != nil {
		return
	}

	var buf [4]byte
	packLen := uint32(len(data))

	fmt.Println("packlen:", packLen)
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

	msg, err = readPackage(conn)
	if err != nil {
		fmt.Println("read package failed, err:", err)
	}

	var registerRes model.RegisterRes
	err = json.Unmarshal([]byte(msg.Data), &registerRes)
	if err != nil {
		fmt.Println("register responce data format not correct !")
	}

	if registerRes.Code == 200 {
		fmt.Println("register successful !")
	}

	if registerRes.Code == 500 {
		fmt.Println("register failed !")
		fmt.Println("err : ", registerRes.Error)
	}
	return
}
