package migateway

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	FIELD_STATUS  = "status"
	FIELD_BATTERY = "battery"
	FIELD_INUSE   = "inuse"
)

type Device struct {
	Sid     string `json:"sid,omitempty"`
	Model   string `json:"model,omitempty"`
	ShortID int    `json:"short_id,omitempty"`
	Data    string `json:"data,omitempty"`

	heartBeatTime int64
	dataMap       map[string]string
}

func (d *Device) RefreshStatus() error {
	return nil
}

func (d *Device) freshHeartTime() {
	d.heartBeatTime = time.Now().Unix()
}

func (d *Device) hasFiled(field string) (found bool) {
	if d.dataMap == nil {
		d.dataMap = make(map[string]string)
		err := json.Unmarshal([]byte(d.Data), &d.dataMap)
		if err != nil {
			LOGGER.Error("parse \"%s\" to map failed: %v", d.Data, err)
			return false
		}
	}

	_, found = d.dataMap[field]
	return
}

func (d *Device) GetData(field string) string {
	if d.hasFiled(field) {
		return d.dataMap[field]
	} else {
		return ""
	}
}

func (d *Device) GetDataAsBool(field string) bool {
	v := d.GetData(field)
	if v == "1" ||
		v == "open" ||
		v == "on" ||
		v == "true" ||
		v == "motion" {
		return true
	} else {
		return false
	}
}

func (d *Device) GetDataAsInt(field string) int {
	v := d.GetData(field)
	n, err := strconv.Atoi(v)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to int failed: %v", v, err)
	}
	return n
}

func (d *Device) GetDataAsUint64(field string) uint64 {
	v := d.GetData(field)
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to uint64 failed: %v", v, err)
	}
	return n
}

func (d *Device) GetDataAsFloat64(field string) float64 {
	v := d.GetData(field)
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to float64 failed: %v", v, err)
	}
	return n
}

func (d *Device) GetDataArray() (array []string) {
	array = make([]string, 0)
	err := json.Unmarshal([]byte(d.Data), &array)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to array failed: %v", d.Data, err)
	}

	return
}
