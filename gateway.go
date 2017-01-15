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
type GateWayDevice struct {
	*Device
	IP   string
	Port string
	RGB  uint32
}

func NewGateWayDevice(dev *Device) *GateWayDevice {
	g := &GateWayDevice{Device: dev}
	g.Set(dev)
	return g
}

func (g *GateWayDevice) Set(dev *Device) {
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

func (gwd *GateWayDevice) ChangeColor(conn *GateWayConn, c color.Color) error {
	data := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	return conn.Control(gwd.Device, data)
}

func (gwd *GateWayDevice) Flashing(conn *GateWayConn, c color.Color) error {
	interval := 500 * time.Millisecond
	err := gwd.flashingOnce(conn, c, interval)
	if err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(interval)
			gwd.flashingOnce(conn, c, interval)
		}
	}()
	return nil
}

func (gwd *GateWayDevice) flashingOnce(conn *GateWayConn, c color.Color, interval time.Duration) error {
	updata := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	downdata := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(COLOR_BLACK)}

	err := conn.Control(gwd.Device, updata)
	if err != nil {
		return err
	}

	time.Sleep(interval)
	err = conn.Control(gwd.Device, downdata)
	if err != nil {
		return err
	}
	return nil
}
