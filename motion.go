package migateway

import (
	"time"
)

const (
	MODEL_MOTION    = "motion"

	FIELD_NO_MOTION = "no_motion"
)

type Motion struct {
	*Device
	State MotionState
}

type MotionState struct {
	Battery    float32
	HasMotion  bool
	LastMotion time.Time
}

type MotionStateChange struct {
	ID   string
	From MotionState
	To   MotionState
}

func (m MotionStateChange) IsChanged() bool {
	return m.From.HasMotion != m.To.HasMotion || m.From.LastMotion != m.To.LastMotion
}

func (m *Motion) GetData() interface{} {
	return m.Data
}

func NewMotion(dev *Device) *Motion {
	return &Motion{
		Device: dev,
		State: MotionState{
			Battery:   dev.GetBatteryLevel(0),
			HasMotion: dev.GetDataAsBool(FIELD_STATUS),
		},
	}
}

func (m *Motion) Set(dev *Device) {
	timestamp := time.Now()
	change := &MotionStateChange{ID: m.Sid, From: m.State, To: m.State}
	if dev.hasField(FIELD_STATUS) {
		m.State.HasMotion = dev.GetDataAsBool(FIELD_STATUS)
		if m.State.HasMotion {
			m.State.LastMotion = timestamp
		}
	} else if dev.hasField(FIELD_NO_MOTION) {
		m.State.HasMotion = false
		nomotionInSeconds := int64(dev.GetDataAsInt(FIELD_NO_MOTION)) * -1
		timestamp.Add(time.Duration(nomotionInSeconds) * time.Second)
		m.State.LastMotion = timestamp
	}

	m.State.Battery = dev.GetBatteryLevel(m.State.Battery)
	change.To = m.State
	if change.IsChanged() || m.shouldPushUpdates() {
		m.Aqara.StateMessages <- change
	}
	if dev.Token != "" {
		m.Token = dev.Token
	}
}
