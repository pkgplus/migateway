package migateway

const (
	MODEL_MOTION = "motion"
)

type MotionDevice struct {
	*DeviceBaseInfo
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
}

func (m *MotionDevice) GetData() interface{} {
	return m.Data
}

func (m *MotionDevice) GetStatus() string {
	return m.Data.Status
}
