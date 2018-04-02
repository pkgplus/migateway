package migateway

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const (
	FIELD_STATUS  = "status"
	FIELD_BATTERY = "battery"
	FIELD_INUSE   = "inuse"
	FIELD_TOKEN   = "token"
)

type Device struct {
	conn    *GateWayConn
	Sid     string `json:"sid,omitempty"`
	Model   string `json:"model,omitempty"`
	ShortID int    `json:"short_id,omitempty"`
	Data    string `json:"data,omitempty"`
	Token   string `json:"token,omitempty"`

	ReportChan    chan interface{} `json:"-"`
	heartBeatTime int64
	dataMap       map[string]interface{}
}

func (d *Device) RefreshStatus() error {
	return nil
}

func (d *Device) freshHeartTime() {
	d.heartBeatTime = time.Now().Unix()
}

func (d *Device) GetHeartTime() int64 {
	return d.heartBeatTime
}

func (d *Device) report(msg interface{}) {
	select {
	case d.ReportChan <- msg:
	default:
	}
}

// func (d *Device) setToken(t string) {
// 	d.Token = t
// }

func (d *Device) waitToken() bool {
	begin := time.Now().Unix()
	for {
		ct := time.Now().Unix()
		if d.conn.token != "" &&
			d.heartBeatTime-ct <= 7200 {
			return true
		} else if ct-begin >= 30 {
			return false
		}
		time.Sleep(time.Second)
	}

	return false
}

func (d *Device) hasField(field string) (found bool) {
	if d.Data == "" {
		return false
	}

	if d.dataMap == nil {
		d.dataMap = make(map[string]interface{})
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
	if d.hasField(field) {
		v, found := d.dataMap[field]
		if found {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Int:
				return fmt.Sprintf("%d", v.(int))
			case reflect.Float64:
				return fmt.Sprintf("%0.f", v.(float64))
			case reflect.String:
				return v.(string)
			default:
				LOGGER.Warn("unknonw %s value type %s", field, reflect.TypeOf(v).Kind().String())
				return ""
			}
		}
	}
	return ""
}

func (d *Device) GetDataAsBool(field string) bool {
	v := d.GetData(field)
	if v == "1" ||
		v == "open" ||
		v == "on" ||
		v == "true" ||
		v == "motion" {
		return true
	}
	return false
}

func (d *Device) GetDataAsInt(field string) int {
	v := d.GetData(field)
	n, err := strconv.Atoi(v)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to int failed: %v", v, err)
	}
	return n
}

func (d *Device) GetDataAsUint32(field string) uint32 {
	v := d.GetData(field)

	n, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		LOGGER.Warn("parse \"%s\" to uint32 failed: %v", v, err)
	}
	return uint32(n)
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
