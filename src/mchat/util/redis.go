package util

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var Pool *redis.Pool

func init() {
	fmt.Println("init redis connect ")
	initRedis("localhost:6379", 16, 1024, time.Second*300)
}

func initRedis(addr string, idleConn, maxConn int, idleTimeout time.Duration) {

	Pool = &redis.Pool{
		MaxIdle:     idleConn,
		MaxActive:   maxConn,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
	return
}

func GetConn() redis.Conn {
	return Pool.Get()
}

func PutConn(conn redis.Conn) {
	conn.Close()
}
