package migateway

const (
	MODEL_PLUG               = "plug"
	FIELD_PLUG_LOADVOLTAGE   = "load_voltage"
	FIELD_PLUG_LOADPOWER     = "load_power"
	FIELD_PLUG_POWERCONSUMED = "power_consumed"
)

type Plug struct {
	*Device
	State PlugState
}

type PlugState struct {
	LoadVoltage   uint32
	PowerConsumed uint32
	LoadPower     float32
	InUse         bool
	IsOn          bool
}

type PlugStateChange struct {
	ID   string
	From PlugState
	To   PlugState
}

func (p PlugStateChange) IsChanged() bool {
	return p.From.InUse != p.To.InUse || p.From.IsOn != p.To.IsOn || p.From.LoadPower != p.To.LoadPower || p.From.LoadVoltage != p.To.LoadVoltage || p.From.PowerConsumed != p.To.PowerConsumed
}

func NewPlug(dev *Device) *Plug {
	p := &Plug{Device: dev}
	p.Set(dev)
	return p
}

func (p *Plug) Set(dev *Device) {
	change := &PlugStateChange{ID: p.Sid, From: p.State, To: p.State}
	if dev.hasField(FIELD_INUSE) {
		p.State.InUse = dev.GetDataAsBool(FIELD_INUSE)
	}
	if dev.hasField(FIELD_STATUS) {
		p.State.IsOn = dev.GetDataAsBool(FIELD_STATUS)
	}
	if dev.hasField(FIELD_PLUG_LOADVOLTAGE) {
		p.State.LoadVoltage = dev.GetDataAsUint32(FIELD_PLUG_LOADVOLTAGE)
	}
	if dev.hasField(FIELD_PLUG_LOADPOWER) {
		p.State.LoadPower = dev.GetDataAsFloat32(FIELD_PLUG_LOADPOWER)
	}
	if dev.hasField(FIELD_PLUG_POWERCONSUMED) {
		p.State.PowerConsumed = dev.GetDataAsUint32(FIELD_PLUG_POWERCONSUMED)
	}
	change.To = p.State
	if change.IsChanged() {
		p.Aqara.StateMessages <- change
	}
	if dev.Token != "" {
		p.Token = dev.Token
	}
}

func (p *Plug) Toggle() error {
	tostate := "on"
	if p.State.IsOn == false {
		tostate = "off"
	}
	data := map[string]interface{}{
		FIELD_STATUS: tostate,
	}
	return p.GatewayConnection.Control(p.Device, data)
}

func (p *Plug) TurnOn() error {
	data := map[string]interface{}{
		FIELD_STATUS: "on",
	}
	return p.GatewayConnection.Control(p.Device, data)
}

func (p *Plug) TurnOff() error {
	data := map[string]interface{}{
		FIELD_STATUS: "off",
	}
	return p.GatewayConnection.Control(p.Device, data)
}
