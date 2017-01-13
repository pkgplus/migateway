package migateway

import (
	"time"
)

type Device interface {
	GetModel() string
	GetData() interface{}
	RefreshStatus() error
	freshHeartTime(int64)
}

type DeviceBaseInfo struct {
	Sid     string `json:"sid"`
	Model   string `json:"model,omitempty"`
	ShortID string `json:"short_id,omitempty"`

	heartBeatTime int64
}

func (d *DeviceBaseInfo) GetModel() string {
	return d.Model
}

func (d *DeviceBaseInfo) RefreshStatus() error {
	return nil
}

func (d *DeviceBaseInfo) freshHeartTime() {
	d.heartBeatTime = time.Now().UnixNano()
}
