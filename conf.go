package migateway

import (
	"time"
)

type Configure struct {
	WhoisTimeOut int
	WhoisRetry   int

	DevListTimeOut int
	DevListRetry   int

	ReadTimeout int
	ReadRetry   int

	AESKey string
}

func (c *Configure) getRetryAndTimeout(req *Request) (int, time.Duration) {
	if req.Cmd == CMD_WHOIS {
		return c.WhoisRetry, time.Duration(c.WhoisTimeOut) * time.Second
	} else if req.Cmd == CMD_DEVLIST {
		return c.DevListRetry, time.Duration(c.DevListTimeOut) * time.Second
	} else {
		return c.ReadRetry, time.Duration(c.ReadTimeout) * time.Second
	}
}

func (c *Configure) SetAESKey(key string) {
	c.AESKey = key
}
