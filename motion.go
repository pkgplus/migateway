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
	NoMotionTime   int
	changeTime     int64
}

func (m *Motion) GetData() interface{} {
	return m.Data
}

func NewMotion(dev *Device) *Motion {
	dev.ReportChan = make(chan interface{}, 1)
	m := &Motion{Device: dev}
	m.Set(dev)
	return m
}

func (m *Motion) Set(dev *Device) {
	last := m.IsMotorial
	ct := time.Now().Unix()

	if dev.hasField(FIELD_STATUS) {
		m.IsMotorial = dev.GetDataAsBool(FIELD_STATUS)

		if m.IsMotorial {
			m.lastMotionTime = ct
		}

	} else if dev.hasField("no_motion") {
		m.IsMotorial = false
		m.NoMotionTime = dev.GetDataAsInt("no_motion")
	}

	if last != m.IsMotorial {
		m.changeTime = ct
	}

	if dev.Token != "" {
		m.Token = dev.Token
	}
}
