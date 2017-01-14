package migateway

const (
	MODEL_MAGNET = "magnet"
)

type Magnet struct {
	*Device
	Opened  bool
	Battery uint64
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
		m.Battery = dev.GetDataAsUint64(FIELD_BATTERY)
	}
}
