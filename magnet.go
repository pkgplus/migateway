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
	dev.ReportChan = make(chan interface{}, 1)
	return &Magnet{
		Device: dev,
		Opened: dev.GetDataAsBool(FIELD_STATUS),
	}
}

func (m *Magnet) Set(dev *Device) {
	if dev.hasFiled(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)
		if status == "open" {
			m.Opened = true
		} else {
			m.Opened = false
		}
	}

	if dev.hasFiled(FIELD_BATTERY) {
		m.Battery = dev.GetDataAsUint32(FIELD_BATTERY)
	}

	if dev.Token != "" {
		m.Token = dev.Token
	}
}
