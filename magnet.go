package migateway

const (
	MODEL_MAGNET = "magnet"

	MAGNET_STATUS_OPEN = "open"
)

type Magnet struct {
	*Device
	State MagnetState
}

type MagnetState struct {
	Opened  bool
	Battery float32
}

type MagnetStateChange struct {
	ID   string
	From MagnetState
	To   MagnetState
}

func (m MagnetStateChange) IsChanged() bool {
	return m.From.Opened != m.To.Opened && m.From.Battery != m.To.Battery
}

func NewMagnet(dev *Device) *Magnet {
	return &Magnet{
		Device: dev,
		State: MagnetState{
			Opened:  dev.GetDataAsBool(FIELD_STATUS),
			Battery: dev.GetBatteryLevel(0),
		},
	}
}

func convertToBatteryPercentage(battery uint32) float32 {
	return float32(battery) / 33.0
}

func (m *Magnet) Set(dev *Device) {
	change := &MagnetStateChange{ID: m.Sid, From: m.State, To: m.State}
	if dev.hasField(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)
		m.State.Opened = status == MAGNET_STATUS_OPEN
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
