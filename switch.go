package migateway

const (
	MODEL_SWITCH = "switch"

	SWITCH_STATUS_CLICK       = "click"
	SWITCH_STATUS_DOUBLECLICK = "double_click"
	SWITCH_STATUS_PRESS       = "long_click_press"
	SWITCH_STATUS_RELEASE     = "long_click_release"
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
	Battery float32
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
	return &Switch{
		Device: dev,
		State:  SwitchState{Battery: dev.GetBatteryLevel(0), Click: NoClick},
	}
}

func (s *Switch) Set(dev *Device) {
	change := &SwitchStateChange{ID: s.Sid, From: s.State, To: s.State}

	s.State.Battery = dev.GetBatteryLevel(s.State.Battery)

	if dev.hasField(FIELD_STATUS) {
		status := dev.GetData(FIELD_STATUS)

		switch status {
		case SWITCH_STATUS_CLICK:
			s.State.Click = SingleClick
		case SWITCH_STATUS_DOUBLECLICK:
			s.State.Click = DoubleClick
		case SWITCH_STATUS_PRESS:
			s.State.Click = LongPressClick
		case SWITCH_STATUS_RELEASE:
			s.State.Click = LongReleaseClick
		default:
			s.State.Click = NoClick
		}
	}

	change.To = s.State
	if change.IsChanged() || s.shouldPushUpdates() {
		s.Aqara.StateMessages <- change
	}

	if dev.Token != "" {
		s.Token = dev.Token
	}
}
