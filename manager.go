package migateway

import (
	"encoding/json"
	"errors"
	"github.com/bingbaba/util/logs"
	"net"
	"time"
)

const (
	MULTICAST_ADDR = "224.0.0.50"
	MULTICAST_PORT = 4321
	SERVER_PORT    = 9898
)

var (
	LOGGER      = logs.GetBlogger()
	EOF    byte = 0

	DefaultConf = &Configure{
		WhoisTimeOut:   3,
		WhoisRetry:     5,
		DevListTimeOut: 3,
		DevListRetry:   5,
		ReadTimeout:    3,
		ReadRetry:      1,
		AESKey:         "pk7clc5c1318qldn",
	}
)

type GateWayManager struct {
	*GateWayDevice
	//IP   string
	Port string

	Motions   map[string]*Motion
	Switchs   map[string]*Switch
	SensorHTs map[string]*SensorHT
	Magnets   map[string]*Magnet
	Plugs     map[string]*Plug

	DiscoveryTime    int64
	FreshDevListTime int64

	SendMsgs   chan []byte
	RecvMsgs   chan []byte
	SendGWMsgs chan []byte
	RecvGWMsgs chan []byte

	Conf *Configure
}

type Configure struct {
	WhoisTimeOut int
	WhoisRetry   int

	DevListTimeOut int
	DevListRetry   int

	ReadTimeout int
	ReadRetry   int

	AESKey string
}

func NewGateWayManager(c *Configure) (m *GateWayManager, err error) {
	if c == nil {
		c = DefaultConf
	}

	m = &GateWayManager{
		Motions:       make(map[string]*Motion),
		Switchs:       make(map[string]*Switch),
		SensorHTs:     make(map[string]*SensorHT),
		Magnets:       make(map[string]*Magnet),
		Plugs:         make(map[string]*Plug),
		DiscoveryTime: time.Now().Unix(),
		SendMsgs:      make(chan []byte),
		RecvMsgs:      make(chan []byte, 100),
		SendGWMsgs:    make(chan []byte),
		RecvGWMsgs:    make(chan []byte, 100),
		Conf:          c,
	}

	//connection
	err = m.initMultiCastConn()
	if err != nil {
		return
	}

	//find the gateway
	m.whois()

	//discovery all device
	err = m.initGateWayConn()
	if err != nil {
		return
	}
	err = m.discovery()

	return
}

func (m *GateWayManager) Control(dev *Device, data map[string]interface{}) error {
	req, err := NewWriteRequest(dev, m.Conf.AESKey, data)
	if err != nil {
		return err
	}
	m.Send(req)
	return nil
}

func (m *GateWayManager) SetAESKey(key string) {
	m.Conf.AESKey = key
}

func (m *GateWayManager) putDevice(dev *Device) (added bool) {
	LOGGER.Info("DEVICESYNC:: %s(%s): %s", dev.Model, dev.Sid, dev.Data)

	added = true
	switch dev.Model {
	case MODEL_GATEWAY:
		m.GateWayDevice.Set(dev)
	case MODEL_MOTION:
		d, found := m.Motions[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			m.Motions[dev.Sid] = NewMotion(dev)
		}
	case MODEL_SWITCH:
		d, found := m.Switchs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			m.Switchs[dev.Sid] = NewSwitch(dev)
		}
	case MODEL_SENSORHT:
		d, found := m.SensorHTs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			m.SensorHTs[dev.Sid] = NewSensorHt(dev)
		}
	case MODEL_MAGNET:
		d, found := m.Magnets[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			m.Magnets[dev.Sid] = NewMagnet(dev)
		}
	case MODEL_PLUG:
		d, found := m.Plugs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			m.Plugs[dev.Sid] = NewPlug(dev)
		}
	default:
		added = false
	}

	return
}

func (m *GateWayManager) whois() {
	//read msg
	iamResp := &IamResp{}
	m.communicate(NewWhoisRequest(), iamResp)

	//gateway infomation
	m.GateWayDevice = NewGateWayDevice(iamResp.Device)
	m.IP = iamResp.IP
	m.Port = iamResp.Port
}

func (m *GateWayManager) discovery() (err error) {
	//get devlist response
	LOGGER.Info("start to discover the device...")
	devListResp := &DeviceListResp{}
	if !m.communicate(NewDevListRequest(), devListResp) {
		return errors.New("show device list error")
	}

	//every device
	for index, sid := range devListResp.getSidArray() {
		dev := m.waitDevice(sid)
		if m.putDevice(dev) {
			LOGGER.Warn("DISCOVERY[%d]: found the device %s(%s): %v", index, dev.Model, dev.Sid, dev.Data)
		} else {
			LOGGER.Warn("DISCOVERY[%d]: unknown model %s device: %v", index, dev.Model, dev)
		}
	}
	m.DiscoveryTime = time.Now().Unix()

	return
}

func (m *GateWayManager) waitDevice(sid string) *Device {
	req := NewReadRequest(sid)
	resp := &DeviceBaseResp{}
	for {
		if m.communicate(req, resp) {
			if resp.Sid == sid {
				break
			} else {
				LOGGER.Info("get a unknown device %s, model = %s, sid = %s", CMD_READ_ACK, resp.Model, resp.Sid)
			}
		}
	}
	return resp.Device
}

func (m *GateWayManager) Send(req *Request) bool {
	if req == nil {
		return false
	} else {
		if req.Cmd == CMD_WHOIS {
			m.multicast(req)
		} else {
			m.sendGW(req)
		}
		return true
	}
}

func (m *GateWayManager) getRetryAndTimeout(req *Request) (int, time.Duration) {
	if req.Cmd == CMD_WHOIS {
		return m.Conf.WhoisRetry, time.Duration(m.Conf.WhoisTimeOut) * time.Second
	} else if req.Cmd == CMD_DEVLIST {
		return m.Conf.DevListRetry, time.Duration(m.Conf.DevListTimeOut) * time.Second
	} else {
		return m.Conf.ReadRetry, time.Duration(m.Conf.ReadTimeout) * time.Second
	}
}

func (m *GateWayManager) getChan(cmd string) chan []byte {
	if cmd == CMD_WHOIS {
		return m.RecvMsgs
	} else {
		return m.RecvGWMsgs
	}
}

func (m *GateWayManager) multicast(req *Request) {
	m.SendMsgs <- toBytes(req)
}

func (m *GateWayManager) sendGW(req *Request) {
	m.SendGWMsgs <- toBytes(req)
}

func (m *GateWayManager) communicate(req *Request, resp Response) bool {
	expectCmd := req.expectCmd()
	if expectCmd == "" {
		LOGGER.Warn("unknown request: %s", string(toBytes(req)))
		return false
	}
	//send message
	m.Send(req)

	retry := 0
	maxRetry, timeout := m.getRetryAndTimeout(req)
	chanName := req.getChanName()

	LOGGER.Info("%s:: wait \"%s\" response...", chanName, expectCmd)
	for {
		select {
		case msg := <-m.getChan(req.Cmd):
			err := json.Unmarshal(msg, resp)
			if err != nil {
				LOGGER.Error("%s:: parse %s error: %v", chanName, string(msg), err)
				continue
			} else if resp.GetCmd() != expectCmd {
				LOGGER.Warn("%s:: wait %s, ingore the msg: %s", chanName, expectCmd, string(msg))
				continue
			} else {
				LOGGER.Info("%s:: recv msg: %s", chanName, string(msg))
				return true
			}
		case <-time.After(time.Second * timeout):
			retry++
			if retry > maxRetry {
				LOGGER.Error("%S:: recv msg TIMEOUT", chanName)
				return false
			} else {
				LOGGER.Error("%S:: send msg retry %d ...", chanName, retry)
				m.Send(req)
			}
		}
	}
	return false
}

func (g *GateWayManager) initGateWayConn() error {
	UDP_ADDR := &net.UDPAddr{
		IP:   net.ParseIP(g.IP),
		Port: SERVER_PORT,
	}
	con, err := net.DialUDP("udp4", nil, UDP_ADDR)
	if err != nil {
		return err
	}

	go func() {
		defer con.Close()
		for {
			msg := <-g.SendGWMsgs
			LOGGER.Info("GATEWAY:: send msg: %s", msg)
			_, werr := con.Write([]byte(msg))
			if werr != nil {
				LOGGER.Error("send error %v", werr)
			}
		}
	}()

	go func() {
		buf := make([]byte, 2048)
		for {
			size, _, err2 := con.ReadFromUDP(buf)
			if err2 != nil {
				panic(err2)
			} else if size > 0 {
				LOGGER.Debug("GATEWAY:: recv msg: %s", string(buf[0:size]))
				g.RecvGWMsgs <- buf[0:size]
			}
		}
	}()

	return nil
}

func (g *GateWayManager) initMultiCastConn() error {
	//listen
	udp_l := &net.UDPAddr{IP: net.ParseIP(MULTICAST_ADDR), Port: SERVER_PORT}
	con, err := net.ListenMulticastUDP("udp4", nil, udp_l)
	if err != nil {
		return err
	}
	LOGGER.Info("listennig %d ...", SERVER_PORT)

	//read
	go func() {
		defer con.Close()

		buf := make([]byte, 2048)
		for {
			size, _, err2 := con.ReadFromUDP(buf)
			if err2 != nil {
				panic(err2)
			} else if size > 0 {
				LOGGER.Debug("MULTICAST:: recv msg: %s", string(buf[0:size]))

				resp := &DeviceBaseResp{}
				json.Unmarshal(buf[0:size], resp)
				// if errTmp != nil {
				// 	LOGGER.Warn("MULTICAST:: parse invalid msg: %s, error:%v", buf[0:size], errTmp)
				// }
				if resp.Cmd == CMD_REPORT || resp.Cmd == CMD_HEARTBEAT {
					g.putDevice(resp.Device)
				} else {
					g.RecvMsgs <- buf[0:size]
				}

			}
		}
	}()

	//write
	MULTI_UDP_ADDR := &net.UDPAddr{
		IP:   net.ParseIP(MULTICAST_ADDR),
		Port: MULTICAST_PORT,
	}
	go func() {
		for {
			msg := <-g.SendMsgs
			wsize, err3 := con.WriteToUDP(msg, MULTI_UDP_ADDR)
			if err3 != nil {
				panic(err3)
			}
			LOGGER.Info("MULTICAST:: send msg: %s, %d bytes!", msg, wsize)
		}
	}()

	return nil
}
