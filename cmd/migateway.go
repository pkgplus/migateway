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
	gwm.SetAESKey("aamfgyc8mbra3jhq")

	time.Sleep(time.Second * 10)

	LOGGER.Info("write...")
	gwd := gwm.GateWayDevice.Device

	err = gwm.Control(gwd, map[string]interface{}{"rgb": 922794751})
	if err != nil {
		panic(err)
	}

	<-make(chan bool)
}
