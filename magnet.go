package migateway

const (
	MODEL_MAGNET = "magnet"
)

type Magnet struct {
	*Device
	Opened  bool
	Battery uint32
}

func NewMagnet(dev *Device) *Magnet {
	return &Magnet{
		Device: dev,
		Opened: dev.GetDataAsBool(FIELD_STATUS),
	}
}

func (m *Magnet) Set(dev *Device) {
	if dev.hasFiled(FIELD_STATUS) {
		m.Opened = dev.GetDataAsBool(FIELD_STATUS)
	}

	if dev.hasFiled(FIELD_BATTERY) {
		m.Battery = dev.GetDataAsUint32(FIELD_BATTERY)
	}

	if dev.Token != "" {
		m.Token = dev.Token
	}
}
