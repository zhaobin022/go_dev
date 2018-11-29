package main

import (
	"fmt"
)

type BalanceMgr struct {
	balanceMap map[string]Balance
}

var mgr BalanceMgr = BalanceMgr{
	balanceMap: make(map[string]Balance),
}

func (b *BalanceMgr) registerBanlance(name string, bal Balance) error {
	mgr.balanceMap[name] = bal

	return nil
}
func RegisterBanlanceMap(name string, bal Balance) {
	mgr.registerBanlance(name, bal)
}

func GetBalance(name string) (b Balance, e error) {
	b, ok := mgr.balanceMap[name]
	if ok != true {
		e = fmt.Errorf("don't has this balance method")
	}
	return
}
