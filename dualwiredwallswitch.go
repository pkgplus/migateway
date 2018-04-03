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
	Channel0On bool
	Channel1On bool
}

func NewDualWiredSwitch(dev *Device) *DualWiredWallSwitch {
	dev.ReportChan = make(chan interface{}, 1)
	return &DualWiredWallSwitch{
		Device:     dev,
		Channel0On: false,
		Channel1On: false,
	}
}

func (s *DualWiredWallSwitch) Set(dev *Device) {

	if dev.hasField(FIELD_CHANNEL0) {
		channel0 := dev.GetDataAsBool(FIELD_CHANNEL0)
		s.Channel0On = channel0
		//LOGGER.Warn("%s", status)
	}
	if dev.hasField(FIELD_CHANNEL1) {
		channel1 := dev.GetDataAsBool(FIELD_CHANNEL1)
		s.Channel1On = channel1
		//LOGGER.Warn("%s", status)
	}
	if dev.Token != "" {
		s.Token = dev.Token
	}
}
