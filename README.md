# migateway
## API DOC
小米网关协议文档参考

[原版](https://github.com/louisZL/lumi-gateway-local-api) 

[修正格式版](https://github.com/xuebing1110/lumi-gateway-local-api)

[绿米网关局域网通讯API V1.0.7](https://cdn.cnbj2.fds.api.mi-img.com/lumiaiot/common/gateway/%E7%BB%BF%E7%B1%B3%E7%BD%91%E5%85%B3%E5%B1%80%E5%9F%9F%E7%BD%91%E9%80%9A%E4%BF%A1%E5%8D%8F%E8%AE%AEV1.0.7_2017.05.25_01.doc)

## FEATURES
* 支持网关自动发现(whois)
* 支持设备自动发现(get_id_list,read)
* 支持设备状态自动更新(report/heartbeat)
* 支持网关彩灯控制，并支持颜色、闪烁等控制(write)
* 支持智能插座的控制
* 服务运行树莓派3上[树莓派安装文档](http://blog.bingbaba.com/post/diy/raspberrypi/)

## Homekit
[migateway for homekit](https://github.com/xuebing1110/gohomekit)

## EXAMPLE
```golang
package main

import (
    "github.com/bingbaba/tool/color"
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
    manager, err := migateway.NewAqaraManager(nil)
    if err != nil {
        panic(err)
    }
    manager.SetAESKey("t7ew6r4y612eml0f")

    gateway := manager.GateWay
    for _, color := range color.COLOR_ALL {
        err = gateway.ChangeColor(color)
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

```

## TODO LIST
* 支持无线开关的控制
* 支持86开关的控制
* 支持场景自动化联动

## Updating...