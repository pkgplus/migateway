package migateway

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	CMD_WHOIS   = `whois`
	CMD_READ    = `read`
	CMD_DEVLIST = `get_id_list`
	CMD_WRITE   = `write`
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
	} else if r.Cmd == CMD_WRITE {
		return CMD_WRITE_ACK
	}
	return ""
}

var (
	iv = []byte{0x17, 0x99, 0x6d, 0x09, 0x3d, 0x28, 0xdd, 0xb3, 0xba, 0x69, 0x5a, 0x2e, 0x6f, 0x58, 0x56, 0x2e}
)

func newWriteRequest(dev *Device, aesKey string, data map[string]interface{}) (*Request, error) {
	//if dev.Token == "" {
	if dev.GatewayConnection.token == "" {
		return nil, errors.New(fmt.Sprintf("the %s(%s) device's token is null!", dev.Model, dev.Sid))
	}

	key_bytes := []byte(aesKey)
	block, err := aes.NewCipher(key_bytes)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, 16)
	//mode.CryptBlocks(ciphertext, []byte(dev.Token))
	mode.CryptBlocks(ciphertext, []byte(dev.GatewayConnection.token))

	data["key"] = fmt.Sprintf("%X", ciphertext)
	bytes, _ := json.Marshal(data)
	dev.Data = string(bytes)
	return &Request{Device: dev, Cmd: CMD_WRITE}, nil
}
