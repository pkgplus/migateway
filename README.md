# migateway
## API DOC
小米网关协议文档参考
[原版](https://github.com/louisZL/lumi-gateway-local-api) 
[修正格式版](https://github.com/xuebing1110/lumi-gateway-local-api)

## FEATURES
* 支持网关自动发现
* 支持设备自动发现
* 支持设备状态自动更新
* 支持网关彩灯控制，并支持颜色控制

## TODO LIST
* 代码优化，调用更简洁
* 增加网关彩灯常用RGB颜色的调用
* 支持场景自动化的配置
* 服务运行树莓派3上[树莓派安装文档](http://blog.bingbaba.com/post/diy/raspberrypi/)

## EXAMPLE
```golang
package main

import (
    "github.com/xuebing1110/migateway"
    "time"
)

func main() {
    //init
    gwm, err := migateway.NewGateWayManager(nil)
    if err != nil {
        panic(err)
    }

    //set the gate aes key
    gwm.SetAESKey("aamfgyc8mbra3jhq")

    //wait 10 secondes
    time.Sleep(time.Second * 10)
    gwd := gwm.GateWayDevice.Device
    err = gwm.Control(gwd, map[string]interface{}{"rgb": 922794751})
    if err != nil {
        panic(err)
    }

    //do something...
}

```

## 持续更新中...