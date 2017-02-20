package migateway

import (
	"time"
)

const (
	MODEL_MOTION = "motion"
)

type Motion struct {
	*Device
	IsMotorial     bool
	lastMotionTime int64
	changeTime     int64
}

func (m *Motion) GetData() interface{} {
	return m.Data
}

func NewMotion(dev *Device) *Motion {
	dev.ReportChan = make(chan bool, 1)
	m := &Motion{Device: dev}
	m.Set(dev)
	return m
}

func (m *Motion) Set(dev *Device) {
	if dev.hasFiled(FIELD_STATUS) {
		last := m.IsMotorial
		ct := time.Now().Unix()
		m.IsMotorial = dev.GetDataAsBool(FIELD_STATUS)

		if m.IsMotorial {
			m.lastMotionTime = ct
		}

		if last != m.IsMotorial {
			m.changeTime = ct
		}
	}

	if dev.Token != "" {
		m.Token = dev.Token
	}
}
