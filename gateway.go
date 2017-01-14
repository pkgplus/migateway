package migateway

const (
	MODEL_GATEWAY     = "gateway"
	FIELD_GATEWAY_RGB = "rgb"
	FIELD_IP          = "ip"
)

//GateWay Status
type GateWayDevice struct {
	*Device
	IP  string
	RGB uint32
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
