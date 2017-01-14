package migateway

const (
	MODEL_GATEWAY     = "gateway"
	FIELD_GATEWAY_RGB = "rgb"
)

//GateWay Status
type GateWayDevice struct {
	*Device
	RGB string
}

func NewGateWayDevice(dev *Device) *GateWayDevice {
	return &GateWayDevice{
		Device: dev,
		RGB:    dev.GetData(FIELD_GATEWAY_RGB),
	}
}

func (g *GateWayDevice) Set(dev *Device) {
	if dev.hasFiled(FIELD_GATEWAY_RGB) {
		g.RGB = dev.GetData(FIELD_GATEWAY_RGB)
	}
}
