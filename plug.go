package migateway

const (
	MODEL_PLUG               = "plug"
	FIELD_PLUG_LOADVOLTAGE   = "load_voltage"
	FIELD_PLUG_LOADPOWER     = "load_power"
	FIELD_PLUG_POWERCONSUMED = "power_consumed"
)

type Plug struct {
	*Device
	InUse         bool
	IsOn          bool
	LoadVoltage   uint64
	LoadPower     uint64
	PowerConsumed uint64
}

func NewPlug(dev *Device) *Plug {
	p := &Plug{Device: dev}
	p.Set(dev)
	return p
}

func (p *Plug) Set(dev *Device) {
	if dev.hasFiled(FIELD_INUSE) {
		p.InUse = dev.GetDataAsBool(FIELD_INUSE)
	}
	if dev.hasFiled(FIELD_STATUS) {
		p.IsOn = dev.GetDataAsBool(FIELD_STATUS)
	}
	if dev.hasFiled(FIELD_PLUG_LOADVOLTAGE) {
		p.LoadVoltage = dev.GetDataAsUint64(FIELD_PLUG_LOADVOLTAGE)
	}
	if dev.hasFiled(FIELD_PLUG_LOADPOWER) {
		p.LoadPower = dev.GetDataAsUint64(FIELD_PLUG_LOADPOWER)
	}
	if dev.hasFiled(FIELD_PLUG_POWERCONSUMED) {
		p.PowerConsumed = dev.GetDataAsUint64(FIELD_PLUG_POWERCONSUMED)
	}
}
