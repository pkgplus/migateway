package migateway

const (
	MODEL_GATEWAY = "gateway"
)

//GateWay Status
type GateWayStatusResp struct {
	*GateWayBaseResp
	Data struct {
		RGB string `json:"rgb"`
	} `json:"data"`
}
