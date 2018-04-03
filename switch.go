package migateway

const (
	MODEL_SWITCH = "switch"

	SWITCH_STATUS_CLICK       = "click"
	SWITCH_STATUS_DOUBLECLICK = "double_click"
)

// ClickType holds the different types of clicks you can perform with a switch
type ClickType int

const (
	// SingleClick is a short single click
	NoClick ClickType = iota
	// SingleClick is a short single click
	SingleClick
	// DoubleClick is a quick double click
	DoubleClick
	// LongPressClick is performed by press and holding the switch
	LongPressClick
	// LongReleaseClick is performed by press, holding and releasing the switch
	LongReleaseClick
)

//GateWay Status
type Switch struct {
	*Device
	Battery int
	Click   ClickType
}

func NewSwitch(dev *Device) *Switch {
	dev.ReportChan = make(chan interface{}, 1)
	return &Switch{
		Device:  dev,
		Battery: dev.GetDataAsInt(FIELD_BATTERY),
	}
}

func (s *Switch) Set(dev *Device) {
	if dev.hasField(FIELD_BATTERY) {
		s.Battery = dev.GetDataAsInt(FIELD_BATTERY)
	}

	if dev.hasField(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)
		s.report(status)

		if status == "click" {
			s.Click = SingleClick
		} else if status == "double_click" {
			s.Click = DoubleClick
		} else if status == "long_click_press" {
			s.Click = LongPressClick
		} else if status == "long_click_release" {
			s.Click = LongReleaseClick
		} else {
			s.Click = NoClick
		}

		//LOGGER.Warn("%s", status)
	}

	if dev.Token != "" {
		s.Token = dev.Token
	}
}
