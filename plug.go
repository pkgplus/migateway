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
	LoadVoltage   uint32
	LoadPower     uint32
	PowerConsumed uint32
}

func NewPlug(dev *Device) *Plug {
	dev.ReportChan = make(chan bool, 1)
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
		p.LoadVoltage = dev.GetDataAsUint32(FIELD_PLUG_LOADVOLTAGE)
	}
	if dev.hasFiled(FIELD_PLUG_LOADPOWER) {
		p.LoadPower = dev.GetDataAsUint32(FIELD_PLUG_LOADPOWER)
	}
	if dev.hasFiled(FIELD_PLUG_POWERCONSUMED) {
		p.PowerConsumed = dev.GetDataAsUint32(FIELD_PLUG_POWERCONSUMED)
	}
	if dev.Token != "" {
		p.Token = dev.Token
	}
}
