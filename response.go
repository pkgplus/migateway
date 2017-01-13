package migateway

const (
	CMD_IAM         = `iam`
	CMD_DEVLIST_ACK = `get_id_list_ack`
	CMD_HEARTBEAT   = `heartbeat`
	CMD_REPORT      = `report`
	CMD_READ_ACK    = `read_ack`

	CMD_READ = `read`
)

type Response interface {
	GetCmd() string
	GetModel() string
}

type DeviceBaseResp struct {
	Device
	Cmd string `json:"cmd"`
}

func (r *DeviceBaseResp) GetCmd() string {
	return r.Cmd
}

//GateWay Discovery: cmd="iam"
type IamResp struct {
	*DeviceBaseResp
	IP   string `json:"ip"`
	Port string `json:"port"`
}

//Device List
type DeviceListResp struct {
	*DeviceBaseResp          //get_id_list_ack
	Data            []string `json:"data"`
}
