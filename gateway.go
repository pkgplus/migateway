package migateway

const (
	MODEL_GATEWAY = "gateway"
)

//GateWay Status
type GateWayDevice struct {
	*Device
	RGB string `json:"rgb"`
}

func (g *GateWayDevice) GetRGB() string {
	return g.RGB
}
