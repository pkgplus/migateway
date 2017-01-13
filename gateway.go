package migateway

const (
	MODEL_GATEWAY = "gateway"
)

//GateWay Status
type GateWayDevice struct {
	*DeviceBaseInfo
	Data struct {
		RGB string `json:"rgb"`
	} `json:"data"`
}

func (g *GateWayDevice) GetData() interface{} {
	return g.Data
}

func (g *GateWayDevice) GetRGB() string {
	return g.Data.RGB
}
