package migateway

const (
	MODEL_SWITCH = "switch"

	SWITCH_STATUS_CLICK       = "click"
	SWITCH_STATUS_DOUBLECLICK = "double_click"
)

//GateWay Status
type Switch struct {
	*Device
	Battery int
}

func NewSwitch(dev *Device) *Switch {
	dev.ReportChan = make(chan interface{}, 1)
	return &Switch{
		Device:  dev,
		Battery: dev.GetDataAsInt(FIELD_BATTERY),
	}
}

func (s *Switch) Set(dev *Device) {
	if dev.hasFiled(FIELD_BATTERY) {
		s.Battery = dev.GetDataAsInt(FIELD_BATTERY)
	}

	if dev.hasFiled(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)
		s.report(status)

		//LOGGER.Warn("%s", status)
	}

	if dev.Token != "" {
		s.Token = dev.Token
	}
}
