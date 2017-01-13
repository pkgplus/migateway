package migateway

import (
	"strconv"
)

const (
	MODEL_SWITCH = "switch"

	SWITCH_STATUS_CLICK       = "click"
	SWITCH_STATUS_DOUBLECLICK = "double_click"
)

//GateWay Status
type Switch struct {
	*Device
	Battery int
	Status  string
}

func NewSwitch(dev *Device) *Switch {
	bstr := dev.GetDataValue("battery")
	bint, _ := strconv.Atoi(bstr)

	return &Switch{
		Device:  dev,
		Battery: bint,
		Status:  dev.GetDataValue("status"),
	}
}
