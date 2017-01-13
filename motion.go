package migateway

const (
	MODEL_MOTION = "motion"
)

type Motion struct {
	*Device
	Status string `json:"status"`
}

func (m *Motion) GetData() interface{} {
	return m.Data
}

func NewMotion(dev *Device) *Motion {
	return &Motion{
		Device: dev,
		Status: dev.GetDataValue("status"),
	}
}
