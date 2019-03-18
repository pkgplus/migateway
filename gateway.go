package migateway

import (
	"errors"
	"image/color"
	"time"

	mcolor "github.com/bingbaba/tool/color"
)

const (
	MODEL_GATEWAY              = "gateway"
	FIELD_GATEWAY_RGB          = "rgb"
	FIELD_GATEWAY_ILLUMINATION = "illumination"
	FIELD_GATEWAY_IP           = "ip"

	FLASHING_WEIGHT_WEAK   = "1"
	FLASHING_WEIGHT_NORMAL = "2"
	FLASHING_WEIGHT_STRONG = "3"
)

//GateWay Status
type GateWay struct {
	*Device
	IP    string
	Port  string
	State GatewayState
}

type GatewayState struct {
	RGB          uint32
	Illumination float64 // Unit: percentage (0-100%)
}
type GatewayStateChange struct {
	ID   string
	From GatewayState
	To   GatewayState
}

func (g GatewayStateChange) IsChanged() bool {
	return g.From.Illumination != g.To.Illumination || g.From.RGB != g.To.RGB
}

func NewGateWay(dev *Device) *GateWay {
	g := &GateWay{Device: dev, State: GatewayState{RGB: 0, Illumination: 0.0}}
	g.Set(dev)
	return g
}

func convertIlluminationToPercentage(illumination int) float64 {
	illumination = (illumination - 300)
	if illumination < 0 {
		illumination = 0
	} else if illumination > 1000 {
		illumination = 1000
	}
	return float64(illumination) / 10.0
}

func (g *GateWay) Set(dev *Device) {
	change := &GatewayStateChange{ID: g.Sid, From: g.State, To: g.State}
	if dev.hasField(FIELD_GATEWAY_IP) {
		g.IP = dev.GetData(FIELD_GATEWAY_IP)
	}
	if dev.hasField(FIELD_GATEWAY_RGB) {
		g.State.RGB = dev.GetDataAsUint32(FIELD_GATEWAY_RGB)
	}
	if dev.hasField(FIELD_GATEWAY_ILLUMINATION) {
		illumination := dev.GetDataAsInt(FIELD_GATEWAY_ILLUMINATION)
		g.State.Illumination = convertIlluminationToPercentage(illumination)
	}

	change.To = g.State
	if change.IsChanged() || g.shouldPushUpdates() {
		g.Aqara.StateMessages <- change
	}

	if dev.Token != "" {
		g.setToken(dev.Token)
	}

	if dev.ShortID > 0 {
		g.ShortID = dev.ShortID
	}
}

func (g *GateWay) setToken(token string) {
	g.Token = token
	g.GatewayConnection.token = token
}

func (gwd *GateWay) ChangeColor(c color.Color) error {
	r, g, b, a := c.RGBA()
	LOGGER.Info("r:%d,g:%d,b:%d,a:%d", r&0xff, g&0xff, b&0xff, a&0xff)
	data := map[string]interface{}{FIELD_GATEWAY_RGB: RGBNumber(c)}
	return gwd.GatewayConnection.Control(gwd.Device, data)
}

func (gwd *GateWay) ChangeBrightness(b int) error {
	if b < 0 || b > 100 {
		return errors.New("Brightness should be between 0~100")
	}
	if b == 0 {
		return gwd.TurnOff()
	}
	rgbNum := uint32(gwd.State.RGB & 0xff)
	rgb := NewRGB(rgbNum)
	rgb.A = uint8(100 - b)
	return gwd.ChangeColor(rgb)
}

func (gwd *GateWay) TurnOff() error {
	return gwd.ChangeColor(mcolor.COLOR_BLACK)
}

func (gwd *GateWay) TurnOn() error {
	rgb := NewRGB(gwd.State.RGB)
	if gwd.State.RGB == 0 {
		rgb = mcolor.COLOR_WHITE
	}
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

	err := gwd.GatewayConnection.Control(gwd.Device, updata)
	if err != nil {
		return err
	}

	time.Sleep(interval)
	err = gwd.GatewayConnection.Control(gwd.Device, downdata)
	if err != nil {
		return err
	}
	return nil
}
