package migateway

import (
	"time"
)

var (
	DefaultConf = &Configure{
		WhoisTimeOut:         3,
		WhoisRetry:           5,
		DevListTimeOut:       3,
		DevListRetry:         5,
		ReadTimeout:          3,
		ReadRetry:            1,
		ReportForwardTimeout: 1,
		ReportAllMessages:    false,
		AESKey:               "",
	}
)

type Configure struct {
	WhoisTimeOut         int
	WhoisRetry           int
	DevListTimeOut       int
	DevListRetry         int
	ReadTimeout          int
	ReadRetry            int
	ReportForwardTimeout int
	ReportAllMessages    bool

	AESKey string
}

func NewConfig() *Configure {
	return &Configure{
		WhoisTimeOut:         3,
		WhoisRetry:           5,
		DevListTimeOut:       3,
		DevListRetry:         5,
		ReadTimeout:          3,
		ReadRetry:            1,
		ReportForwardTimeout: 1,
		ReportAllMessages:    false,
		AESKey:               "",
	}
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
