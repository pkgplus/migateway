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
	return &Switch{
		Device:  dev,
		Battery: dev.GetDataAsInt(FIELD_BATTERY),
	}
}

func (s *Switch) Set(dev *Device) {
	if s.hasFiled(FIELD_BATTERY) {
		s.Battery = s.GetDataAsInt(FIELD_BATTERY)
	}
	if dev.Token != "" {
		s.Token = dev.Token
	}
}

func (s *Switch) isClick(dev *Device) bool {
	if dev.GetData(FIELD_STATUS) == SWITCH_STATUS_CLICK {
		return true
	} else {
		return false
	}
}

func (s *Switch) isDoubleClick(dev *Device) bool {
	if dev.GetData(FIELD_STATUS) == SWITCH_STATUS_DOUBLECLICK {
		return true
	} else {
		return false
	}
}
