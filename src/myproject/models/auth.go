package models

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func GetRandomSalt() string {
	return GetRandomString(8)
}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GenCryptPassStr(password string) (encryPass string) {
	t := sha1.New()
	io.WriteString(t, password)
	encryPass = fmt.Sprintf("%x", t.Sum(nil))
	return
}

func GetEncryPass(pass string) (encrypass string) {
	salt := GetRandomSalt()
	passFormat := fmt.Sprintf("%s|%s", salt, pass)
	encryPass := GenCryptPassStr(passFormat)
	encrypass = fmt.Sprintf("%s|%s", salt, encryPass)
	return
}
