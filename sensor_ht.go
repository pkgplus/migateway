package migateway

const (
	MODEL_SENSORHT = "sensor_ht"

	FIELD_SENSORHT_TEMPERATURE = "temperature"
	FIELD_SENSORHT_HUMIDITY    = "humidity"
)

type SensorHT struct {
	*Device
	Temperature float64
	Humidity    float64
}

func NewSensorHt(dev *Device) *SensorHT {
	return &SensorHT{
		Device:      dev,
		Temperature: dev.GetDataAsFloat64(FIELD_SENSORHT_TEMPERATURE) / 100,
		Humidity:    dev.GetDataAsFloat64(FIELD_SENSORHT_HUMIDITY) / 100,
	}
}

func (s *SensorHT) Set(dev *Device) {
	if dev.hasFiled(FIELD_SENSORHT_TEMPERATURE) {
		s.Temperature = dev.GetDataAsFloat64(FIELD_SENSORHT_TEMPERATURE) / 100
	}

	if dev.hasFiled(FIELD_SENSORHT_HUMIDITY) {
		s.Humidity = dev.GetDataAsFloat64(FIELD_SENSORHT_HUMIDITY) / 100
	}
	if dev.Token != "" {
		s.Token = dev.Token
	}
}
