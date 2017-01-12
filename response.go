package migateway

type GateWayResp interface {
	GetCmd() string
}

type CommonResp struct {
	Cmd string `json:"cmd"`
	Sid string `json:"sid"`
}

type DeviceBaseResp struct {
	*CommonResp
	Model   string `json:"model"`
	ShortID string `json:"short_id"`
}

//GateWay Discovery
type GateWayDiscResp struct {
	*DeviceBaseResp        //cmd = "iam"
	IP              string `json:"ip"`
	Port            string `json:"port"`
}

//Device List
type GateWayDevListResp struct {
	*CommonResp          //get_id_list_ack
	Data        []string `json:"data"`
}

func (r *CommonResp) GetCmd() string {
	return r.Cmd
}
