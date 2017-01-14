package migateway

import (
	"encoding/json"
)

const (
	CMD_WHOIS   = `whois`
	CMD_READ    = `read`
	CMD_DEVLIST = `get_id_list`
)

type Request struct {
	*Device
	Cmd string `json:"cmd"`
}

func NewReadRequest(sid string) *Request {
	return &Request{Cmd: CMD_READ, Device: &Device{Sid: sid}}
}

func NewWhoisRequest() *Request {
	return &Request{Cmd: CMD_WHOIS}
}

func NewDevListRequest() *Request {
	return &Request{Cmd: CMD_DEVLIST}
}

func toBytes(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

func (r *Request) getChanName() string {
	if r.Cmd == CMD_WHOIS {
		return "MULTICAST"
	} else {
		return "GATEWAY"
	}
}

func (r *Request) expectCmd() string {
	if r.Cmd == CMD_WHOIS {
		return CMD_IAM
	} else if r.Cmd == CMD_READ {
		return CMD_READ_ACK
	} else if r.Cmd == CMD_DEVLIST {
		return CMD_DEVLIST_ACK
	}
	return ""
}
