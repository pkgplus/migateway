package migateway

import (
	"errors"
	mcolor "github.com/bingbaba/tool/color"
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
	IP       string
	Port     string
	lastRGB  uint32
	RGB      uint32
	callBack func(gw *GateWay) error
}

func NewGateWay(dev *Device) *GateWay {
	dev.ReportChan = make(chan bool, 1)
	g := &GateWay{Device: dev}
	g.Set(dev)
	return g
}

func (g *GateWay) Set(dev *Device) {
	if dev.hasFiled(FIELD_IP) {
		g.IP = dev.GetData(FIELD_IP)
	}
	if dev.hasFiled(FIELD_GATEWAY_RGB) {
		if g.RGB != RGBNumber(mcolor.COLOR_BLACK) && g.RGB != 0 {
			LOGGER.Info("Save Last RGB:%d", g.RGB)
			g.lastRGB = g.RGB
		}
		g.RGB = dev.GetDataAsUint32(FIELD_GATEWAY_RGB)
	}
	if dev.Token != "" {
		g.Token = dev.Token
	}
	if dev.ShortID > 0 {
		g.ShortID = dev.ShortID
	}
	if g.callBack != nil {
		err := g.callBack(g)
		if err != nil {
			LOGGER.Error("exec callback error:%v", err)
		}
	}
}

func (gwd *GateWay) RegisterCb(cb func(gw *GateWay) error) {
	gwd.callBack = cb
}

func (gwd *GateWay) ChangeColor(c color.Color) error {
	data := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	return gwd.conn.Control(gwd.Device, data)
}

func (gwd *GateWay) ChangeBrightness(b int) error {
	if b < 0 || b > 100 {
		return errors.New("Brightness should be 0~100")
	}

	if b == 0 {
		return gwd.TurnOff()
	} else {
		rgb := NewRGB(gwd.RGB)
		rgb.Brightness(uint8(b))
		return gwd.ChangeColor(rgb)
	}
}

func (gwd *GateWay) TurnOff() error {
	return gwd.ChangeColor(mcolor.COLOR_BLACK)
}

func (gwd *GateWay) TurnOn() error {
	if gwd.RGB > 0 {
		gwd.lastRGB = gwd.RGB
	} else if gwd.lastRGB == 0 {
		gwd.lastRGB = RGBNumber(mcolor.COLOR_WHITE)
		LOGGER.Warn("Use Default RGB:%d", gwd.lastRGB)
	}

	rgb := NewRGB(gwd.lastRGB)
	r, g, b, a := rgb.RGBA_8bit()
	LOGGER.Warn("Change Color %d,%d,%d,%d", r, g, b, a)
	return gwd.ChangeColor(rgb)
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
	downdata := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(mcolor.COLOR_BLACK)}

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
