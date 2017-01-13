package migateway

import (
	"strconv"
)

const (
	MODEL_SENSORHT = "sensor_hb"
)

type SensorHT struct {
	*Device
	Temperature float64
	Humidity    float64
}

func NewSensorHt(dev *Device) *SensorHT {
	tStr := dev.GetDataValue("temperature")
	t, _ := strconv.ParseFloat(tStr, 64)
	hStr := dev.GetDataValue("humidity")
	h, _ := strconv.ParseFloat(hStr, 64)

	return &SensorHT{
		Device:      dev,
		Temperature: t / 100,
		Humidity:    h / 100,
	}
}
