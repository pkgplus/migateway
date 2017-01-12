package migateway

import (
	"encoding/json"
)

type GateWayRequest struct {
	Cmd   string      `json:"cmd"`
	Sid   string      `json:"sid"`
	Model string      `json:"model,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func NewGateWayReadRequest(sid string) string {
	bytes, _ := json.Marshal(&GateWayRequest{Cmd: CMD_READ, Sid: sid})
	return string(bytes)
}
