package migateway

import (
	"encoding/json"
)

type GwRequest struct {
	Cmd string `json:"cmd"`
	Sid string `json:"sid"`
}

func NewReadRequest(sid string) string {
	bytes, _ := json.Marshal(&GwRequest{Cmd: CMD_READ, Sid: sid})
	return string(bytes)
}
