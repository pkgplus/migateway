package main

import (
	"github.com/xuebing1110/migateway"
)

func main() {
	_, err := migateway.NewGateWayManager()
	if err != nil {
		panic(err)
	}
}
