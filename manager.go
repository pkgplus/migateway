package migateway

import (
	"errors"
	"github.com/bingbaba/util/logs"
	"time"
)

var (
	LOGGER      = logs.GetBlogger()
	EOF    byte = 0
)

type AqaraManager struct {
	reportChan chan *Device

	GateWay   *GateWay
	Motions   map[string]*Motion
	Switchs   map[string]*Switch
	SensorHTs map[string]*SensorHT
	Magnets   map[string]*Magnet
	Plugs     map[string]*Plug

	DiscoveryTime    int64
	FreshDevListTime int64
}

func NewAqaraManager(c *Configure) (m *AqaraManager, err error) {
	if c == nil {
		c = DefaultConf
	}
	conn := NewConn(c)

	//connection
	err = conn.initMultiCast()
	if err != nil {
		return
	}

	//AqaraManager
	m = &AqaraManager{
		Motions:       make(map[string]*Motion),
		Switchs:       make(map[string]*Switch),
		SensorHTs:     make(map[string]*SensorHT),
		Magnets:       make(map[string]*Magnet),
		Plugs:         make(map[string]*Plug),
		DiscoveryTime: time.Now().Unix(),
	}

	//find gateway
	m.whois(conn)

	//show device list
	gw_ip := m.GateWay.IP
	err = conn.initGateWay(gw_ip)
	if err != nil {
		return
	}

	err = m.discovery()

	//report or heartbeat message
	go func() {
		for {
			m.putDevice(<-conn.devMsgs)
		}
	}()

	return
}

func (m *AqaraManager) putDevice(dev *Device) (added bool) {
	LOGGER.Info("DEVICESYNC:: %s(%s): %s", dev.Model, dev.Sid, dev.Data)
	gateway := m.GateWay

	added = true
	switch dev.Model {
	case MODEL_GATEWAY:
		gateway.Set(dev)
	case MODEL_MOTION:
		d, found := m.Motions[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.Motions[dev.Sid] = NewMotion(dev)
		}
	case MODEL_SWITCH:
		d, found := m.Switchs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.Switchs[dev.Sid] = NewSwitch(dev)
		}
	case MODEL_SENSORHT:
		d, found := m.SensorHTs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.SensorHTs[dev.Sid] = NewSensorHt(dev)
		}
	case MODEL_MAGNET:
		d, found := m.Magnets[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.Magnets[dev.Sid] = NewMagnet(dev)
		}
	case MODEL_PLUG:
		d, found := m.Plugs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.conn = gateway.conn
			m.Plugs[dev.Sid] = NewPlug(dev)
		}
	default:
		added = false
		LOGGER.Warn("DEVICESYNC:: unknown model is %s", dev.Model)
	}

	return
}

func (m *AqaraManager) whois(conn *GateWayConn) {
	//read msg
	iamResp := &IamResp{}
	conn.communicate(NewWhoisRequest(), iamResp)

	//gateway infomation
	dev := NewGateWay(iamResp.Device)
	dev.IP = iamResp.IP
	dev.Port = iamResp.Port
	dev.conn = conn

	m.GateWay = dev
}

func (m *AqaraManager) discovery() (err error) {
	gateway := m.GateWay
	conn := gateway.conn

	//get devlist response
	LOGGER.Info("start to discover the device...")
	devListResp := &DeviceListResp{}
	if !conn.communicate(NewDevListRequest(), devListResp) {
		return errors.New("show device list error")
	}
	gateway.setToken(devListResp.Token)

	//every device
	for index, sid := range devListResp.getSidArray() {
		dev := conn.waitDevice(sid)
		dev.Token = devListResp.Token
		if m.putDevice(dev) {
			LOGGER.Warn("DISCOVERY[%d]: found the device %s(%s): %v", index, dev.Model, dev.Sid, dev.Data)
		} else {
			LOGGER.Warn("DISCOVERY[%d]: unknown model %s device: %v", index, dev.Model, dev)
		}
	}
	m.DiscoveryTime = time.Now().Unix()

	return
}

func (m *AqaraManager) SetAESKey(key string) {
	m.GateWay.conn.SetAESKey(key)
}
