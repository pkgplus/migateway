package main

import (
	"time"

	"github.com/bingbaba/tool/color"
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/migateway"
)

var (
	LOGGER = logs.GetBlogger()
)

func init() {
	logs.SetDebug(true)
}

func main() {
	manager := migateway.NewAqaraManager()
	cfg := migateway.NewConfig()
	cfg.SetAESKey("09F859F7B23A46BE")
	cfg.ReportAllMessages = true
	err := manager.Start(cfg)
	if err != nil {
		panic(err)
	}

	time.Sleep(3 * time.Second)
	manager.Stop()
	time.Sleep(1 * time.Second)

	manager = migateway.NewAqaraManager()
	err = manager.Start(cfg)
	if err != nil {
		panic(err)
	}

	go func() {
		for msg := range manager.StateMessages {
			LOGGER.Info("%+v", msg)
		}
	}()

	gateway := manager.GateWay
	for _, c := range color.COLOR_ALL {
		err := gateway.ChangeColor(c)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}

	err = gateway.Flashing(color.COLOR_RED)
	if err != nil {
		panic(err)
	}

	//do something...
	time.Sleep(10 * time.Second)
}
