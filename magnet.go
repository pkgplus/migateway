package migateway

const (
	MODEL_MAGNET = "magnet"
)

type Magnet struct {
	*Device
	State MagnetState
}

type MagnetState struct {
	Opened  bool
	Battery float64
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
		State:  MagnetState{Opened: dev.GetDataAsBool(FIELD_STATUS), Battery: 100},
	}
}

func convertToBatteryPercentage(battery uint32) float64 {
	return float64(battery)
}

func (m *Magnet) Set(dev *Device) {
	change := &MagnetStateChange{ID: m.Sid, From: m.State, To: m.State}
	if dev.hasField(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)
		if status == "open" {
			m.State.Opened = true
		} else {
			m.State.Opened = false
		}
	}

	if dev.hasField(FIELD_BATTERY) {
		m.State.Battery = convertToBatteryPercentage(dev.GetDataAsUint32(FIELD_BATTERY))
	}

	change.To = m.State
	if change.IsChanged() {
		m.Aqara.StateMessages <- change
	}

	if dev.Token != "" {
		m.Token = dev.Token
	}
}
