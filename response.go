package migateway

const (
	CMD_IAM         = `iam`
	CMD_DEVLIST_ACK = `get_id_list_ack`
	CMD_HEARTBEAT   = `heartbeat`
	CMD_REPORT      = `report`
	CMD_READ_ACK    = `read_ack`
)

type Response interface {
	GetCmd() string
}

type DeviceBaseResp struct {
	*Device
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
	*DeviceBaseResp //get_id_list_ack
	sidArray        []string
}

func (d *DeviceListResp) getSidArray() []string {
	if d.sidArray == nil {
		d.sidArray = d.GetDataArray()
	}
	return d.sidArray
}
