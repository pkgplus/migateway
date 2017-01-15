package main

import (
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/migateway"
	"time"
)

var (
	LOGGER = logs.GetBlogger()
)

func init() {
	logs.SetDebug(true)
}

func main() {
	gwm, err := migateway.NewGateWayManager(nil)
	if err != nil {
		panic(err)
	}

	conn := gwm.GateWayConn
	conn.SetAESKey("t7ew6r4y612eml0f")

	for _, color := range migateway.COLOR_ALL {
		err = gwm.ChangeColor(conn, color)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	err = gwm.Flashing(conn, migateway.COLOR_RED)
	if err != nil {
		panic(err)
	}

	//do something...
	time.Sleep(10 * time.Second)
}
