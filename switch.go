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

func (c ClickType) String() string {
	switch c {
	case SingleClick:
		return "single"
	case DoubleClick:
		return "double"
	case LongPressClick:
		return "longpress"
	case LongReleaseClick:
		return "longrelease"
	}
	return "none"
}

//GateWay Status
type Switch struct {
	*Device
	State SwitchState
}

type SwitchState struct {
	Battery float64
	Click   ClickType
}

type SwitchStateChange struct {
	ID   string
	From SwitchState
	To   SwitchState
}

func (s SwitchStateChange) IsChanged() bool {
	return s.From.Battery != s.To.Battery || s.From.Click != s.To.Click
}

func NewSwitch(dev *Device) *Switch {
	battery := convertToBatteryPercentage(dev.GetDataAsUint32(FIELD_BATTERY))
	return &Switch{
		Device: dev,
		State:  SwitchState{Battery: battery, Click: NoClick},
	}
}

func (s *Switch) Set(dev *Device) {
	if dev.hasField(FIELD_BATTERY) {
		battery := convertToBatteryPercentage(dev.GetDataAsUint32(FIELD_BATTERY))
		s.State.Battery = battery
	}

	change := SwitchStateChange{ID: s.Sid, From: s.State, To: s.State}
	if dev.hasField(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)

		if status == "click" {
			s.State.Click = SingleClick
		} else if status == "double_click" {
			s.State.Click = DoubleClick
		} else if status == "long_click_press" {
			s.State.Click = LongPressClick
		} else if status == "long_click_release" {
			s.State.Click = LongReleaseClick
		} else {
			s.State.Click = NoClick
		}
		//LOGGER.Warn("%s", status)
	}

	change.To = s.State
	if change.IsChanged() {
		s.Gateway.StateChanges <- change
	}

	if dev.Token != "" {
		s.Token = dev.Token
	}
}
