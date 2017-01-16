package migateway

import (
	"image/color"
	"time"
)

const (
	MODEL_GATEWAY     = "gateway"
	FIELD_GATEWAY_RGB = "rgb"
	FIELD_IP          = "ip"

	FLASHING_WEIGHT_WEAK   = "1"
	FLASHING_WEIGHT_NORMAL = "2"
	FLASHING_WEIGHT_STRONG = "3"
)

//GateWay Status
type GateWay struct {
	*Device
	IP   string
	Port string
	RGB  uint32
}

func NewGateWay(dev *Device) *GateWay {
	g := &GateWay{Device: dev}
	g.Set(dev)
	return g
}

func (g *GateWay) Set(dev *Device) {
	if dev.hasFiled(FIELD_IP) {
		g.IP = dev.GetData(FIELD_IP)
	}
	if dev.hasFiled(FIELD_GATEWAY_RGB) {
		g.RGB = dev.GetDataAsUint32(FIELD_GATEWAY_RGB)
	}
	if dev.Token != "" {
		g.Token = dev.Token
	}
	if dev.ShortID > 0 {
		g.ShortID = dev.ShortID
	}
}

func (gwd *GateWay) ChangeColor(c color.Color) error {
	data := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	return gwd.conn.Control(gwd.Device, data)
}

func (gwd *GateWay) Flashing(c color.Color) error {
	interval := 500 * time.Millisecond
	err := gwd.flashingOnce(c, interval)
	if err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(interval)
			gwd.flashingOnce(c, interval)
		}
	}()
	return nil
}

func (gwd *GateWay) flashingOnce(c color.Color, interval time.Duration) error {
	updata := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	downdata := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(COLOR_BLACK)}

	err := gwd.conn.Control(gwd.Device, updata)
	if err != nil {
		return err
	}

	time.Sleep(interval)
	err = gwd.conn.Control(gwd.Device, downdata)
	if err != nil {
		return err
	}
	return nil
}
