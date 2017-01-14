package main

import (
	"github.com/xuebing1110/migateway"
)

func main() {
	_, err := migateway.NewGateWayManager(nil)
	if err != nil {
		panic(err)
	}
}
