package main

import (
	"fmt"
	"hash/crc64"
	"math/rand"
)

func init() {
	var rb RandomBalance
	var rbb RobinBalance
	RegisterBanlanceMap("random", &rb)
	RegisterBanlanceMap("robin", &rbb)
}

type Balance interface {
	DoBalance([]Instance, ...string) (*Instance, error)
}

type RandomBalance struct{}

func (r *RandomBalance) DoBalance(ins []Instance, args ...string) (in *Instance, e error) {
	n := rand.Intn(len(ins))
	in = &ins[n]
	return
}

type RobinBalance struct {
	count int
}

func (r *RobinBalance) DoBalance(ins []Instance, args ...string) (in *Instance, e error) {
	if r.count >= len(ins) {
		r.count = 0
	}
	in = &ins[r.count]
	r.count += 1
	return
}

type HashBalance struct {
	count int
}

func (r *HashBalance) DoBalance(ins []Instance, args ...string) (in *Instance, e error) {
	if len(args) == 0 {
		e = fmt.Errorf("must input the balance key")
		return
	}

	if len(ins) == 0 {
		e = fmt.Errorf("No Instance can use !")
		return
	}
	key := args[0]
	// crcTable := crc64.MakeTable(crc64.ECMA)
	hashVal := crc64.Checksum([]byte(key), nil)
	index := int(hashVal) % len(ins)
	in = &ins[index]
	return
}
