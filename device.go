package migateway

const (
	DEV_MODEL_GATEWAY = "gateway"
)

type Device interface {
	RefreshStatus() error
}

type DeviceBase struct {
	Sid           string
	IP            string
	Model         string
	HeartBeatTime int64
}

type GateWay struct {
	*DeviceBase
	RGB string
}
