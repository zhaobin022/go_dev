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

var UserId int

func login(conn net.Conn) (err error) {
	var msg model.Message
	msg.Cmd = util.UserLogin

	var loginCmd model.LoginCmd
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
		reader = bufio.NewReader(os.Stdin)
		pass, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input password error : ", err)
			continue
		}
		pass = strings.Trim(pass, "\r\n")
		loginCmd.Passwd = pass
		loginCmd.Id = id
		break
	}

	data, err := json.Marshal(loginCmd)
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

	msg, err = readPackage(conn)
	if err != nil {
		fmt.Println("read package failed, err:", err)
	}

	var loginRes model.LoginCmdRes
	err = json.Unmarshal([]byte(msg.Data), &loginRes)
	if err != nil {
		fmt.Println("receive json format not corrent !", err)
	}

	if loginRes.Code == 200 {
		fmt.Println("login success full !")
		UserId = loginCmd.Id
		loginProcess(conn)
	}

	if loginRes.Code == 500 {
		fmt.Println("login fialed ")
		fmt.Println("err : ", loginRes.Error)
	}

	return
}

func processServerNotify(conn net.Conn) {
	for {
		msg, err := readPackage(conn)
		if err != nil {
			fmt.Println("receive server notify err : ", err)
		}
		switch msg.Cmd {
		case util.UserOnlineNotify:
			var user model.User
			err = json.Unmarshal([]byte(msg.Data), &user)
			if err != nil {
				fmt.Println("unfotmat user json in notify case ")
				break
			}
			fmt.Println("user : ", strings.Trim(user.Nick, "\r\n"), " is online !")

		case util.DisplayOnlineUserCmd:
			displayUserList(msg.Data)
		case util.TalkToallCmd:
			displayTalkMsg(conn, msg)
		}

	}

}

func loginProcess(conn net.Conn) {
	go processServerNotify(conn)
	for {
		afterLoginMemu()
		reader := bufio.NewReader(os.Stdin)
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("input option error : ", err)
			continue
		}
		option = strings.Trim(option, "\r\n")
		switch option {
		case "1":
			displayOnlineuser(conn)
		case "2":
			talkToall(conn)
		}
	}

}
