package main

import (
	. "collectlog/conf"
	. "collectlog/kafka"
	. "collectlog/tailf"
)

func main() {
	InitConf()
	InitTail()
	InitKafka()
}
