package main

import (
	"encoding/json"
	"fmt"
	. "mchat/model"
	. "mchat/util"
	"strings"
)

// var (
// 	CliMgr =
// )

type ClientMgr struct {
	allClientMap map[int]Client
}

func (cmr *ClientMgr) AddClientToMap(client Client) (err error) {

	if _, ok := cmr.allClientMap[client.user.UserId]; ok == false {
		cmr.allClientMap[client.user.UserId] = client
	} else {
		err = fmt.Errorf("user already in the map")
	}
	return
}

func (cmr *ClientMgr) NotifyOnline(client Client) {
	user := client.user

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("online user json format error : ", err)
		return
	}

	for k, v := range cmr.allClientMap {
		if k == client.user.UserId {
			continue
		}
		var msg = &Message{Cmd: UserOnlineNotify, Data: string(data)}

		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("online user json format error : ", err)
		}
		fmt.Println("before send onlie notify to ", v.user.Nick)
		err = v.writePackage(data)
		if err != nil {
			fmt.Println("send onlie notify to ", v.user.Nick, " failed")
		}

	}
}

func (cmr *ClientMgr) DisplayOnlineUser(client *Client) (err error) {
	var user_map map[int]string = make(map[int]string)
	for k, v := range cmr.allClientMap {
		user_map[k] = strings.Trim(v.user.Nick, "\r\n")
	}

	data, err := json.Marshal(user_map)
	if err != nil {
		fmt.Println("fotmat user map to json failed !")
		return
	}

	var msg Message

	msg.Cmd = DisplayOnlineUserCmd
	msg.Data = string(data)
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("for online user msg failed !")
		return
	}
	err = client.writePackage(data)
	if err != nil {
		fmt.Println("write online user data failed !")
		return
	}
	return
}

func (cmr *ClientMgr) SendMsgToAll(c *Client, msg Message) (err error) {
	var data_list []string
	var data []byte
	data_list = append(data_list, c.user.Nick)
	data_list = append(data_list, msg.Data)
	data, err = json.Marshal(data_list)
	if err != nil {
		fmt.Println("fotmat ")
	}
	msg.Data = string(data)
	data, err = json.Marshal(msg)

	if err != nil {
		fmt.Println("format msg data error : ", err)
		return
	}

	for k, v := range cmr.allClientMap {
		if k == c.user.UserId {
			continue
		}
		err = v.writePackage(data)
		if err != nil {
			fmt.Printf("send data to conn failed !")
			return
		}
	}

	return
}
