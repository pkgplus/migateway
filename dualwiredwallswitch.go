package migateway

const (
	MODEL_DUALWIREDSWITCH = "ctrl_neutral2"
	FIELD_CHANNEL0        = "channel_0"
	FIELD_CHANNEL1        = "channel_1"
)

// DualWiredWallSwitch is a wired wall switch with 2 switches to control to power-lines.
// The switches can be controlled by the gateway using the zigbee protocol.
type DualWiredWallSwitch struct {
	*Device
	State DualWiredWallSwitchState
}

type DualWiredWallSwitchState struct {
	Channel0On bool
	Channel1On bool
}

type DualWiredWallSwitchStateChange struct {
	ID   string
	From DualWiredWallSwitchState
	To   DualWiredWallSwitchState
}

func (d DualWiredWallSwitchStateChange) IsChanged() bool {
	return d.From.Channel0On != d.To.Channel0On || d.From.Channel1On != d.To.Channel1On
}

func NewDualWiredSwitch(dev *Device) *DualWiredWallSwitch {
	return &DualWiredWallSwitch{
		Device: dev,
		State:  DualWiredWallSwitchState{Channel0On: false, Channel1On: false},
	}
}

func (s *DualWiredWallSwitch) Set(dev *Device) {
	change := &DualWiredWallSwitchStateChange{ID: s.Sid, From: s.State, To: s.State}
	if dev.hasField(FIELD_CHANNEL0) {
		channel0 := dev.GetDataAsBool(FIELD_CHANNEL0)
		s.State.Channel0On = channel0
		change.To = s.State
		//LOGGER.Warn("%s", status)
	}
	if dev.hasField(FIELD_CHANNEL1) {
		channel1 := dev.GetDataAsBool(FIELD_CHANNEL1)
		s.State.Channel1On = channel1
		change.To = s.State
		//LOGGER.Warn("%s", status)
	}

	if change.IsChanged() || s.shouldPushUpdates() {
		s.Aqara.StateMessages <- change
	}

	if dev.Token != "" {
		s.Token = dev.Token
	}
}

func (s *DualWiredWallSwitch) TurnOn() error {
	data := map[string]interface{}{
		"channel_0": "on",
		"channel_1": "on",
	}
	return s.GatewayConnection.Control(s.Device, data)
}
func (s *DualWiredWallSwitch) TurnOnChannel0() error {
	data := map[string]interface{}{
		"channel_0": "on",
	}
	return s.GatewayConnection.Control(s.Device, data)
}
func (s *DualWiredWallSwitch) TurnOnChannel1() error {
	data := map[string]interface{}{
		"channel_1": "on",
	}
	return s.GatewayConnection.Control(s.Device, data)
}

func (s *DualWiredWallSwitch) TurnOff() error {
	data := map[string]interface{}{
		"channel_0": "off",
		"channel_1": "off",
	}
	return s.GatewayConnection.Control(s.Device, data)
}
func (s *DualWiredWallSwitch) TurnOffChannel0() error {
	data := map[string]interface{}{
		"channel_0": "off",
	}
	return s.GatewayConnection.Control(s.Device, data)
}
func (s *DualWiredWallSwitch) TurnOffChannel1() error {
	data := map[string]interface{}{
		"channel_1": "off",
	}
	return s.GatewayConnection.Control(s.Device, data)
}
