package migateway

import (
	"encoding/json"
	"time"
)

// type Device interface {
// 	GetSid() string
// 	GetModel() string
// 	GetShortID() string
// 	GetData() interface{}
// 	RefreshStatus() error
// 	freshHeartTime()
// }

type Device struct {
	Sid     string `json:"sid"`
	Model   string `json:"model,omitempty"`
	ShortID int    `json:"short_id,omitempty"`
	Data    string

	heartBeatTime int64
	dataMap       map[string]string
}

func (d *Device) RefreshStatus() error {
	return nil
}

func (d *Device) freshHeartTime() {
	d.heartBeatTime = time.Now().UnixNano()
}

func (d *Device) GetDataValue(field string) string {
	if d.dataMap == nil {
		d.dataMap = make(map[string]string)
		err := json.Unmarshal([]byte(d.Data), &d.dataMap)
		if err != nil {
			LOGGER.Error("parse %s error: %v", d.Data, err)
			return ""
		}
	}
	return d.dataMap[field]
}

func (d *Device) GetDataArray() (array []string) {
	array = make([]string, 0)
	err := json.Unmarshal([]byte(d.Data), &array)
	if err != nil {
		LOGGER.Error("parse %s error: %v", d.Data, err)
	}

	return
}
