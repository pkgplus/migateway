package migateway

import (
	"errors"
	"github.com/bingbaba/util/logs"
	"time"
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
	*GateWayConn
	*GateWayDevice
	//IP   string
	Port string

	reportChan chan *Device
	Motions    map[string]*Motion
	Switchs    map[string]*Switch
	SensorHTs  map[string]*SensorHT
	Magnets    map[string]*Magnet
	Plugs      map[string]*Plug

	DiscoveryTime    int64
	FreshDevListTime int64
}

func NewGateWayManager(c *Configure) (m *GateWayManager, err error) {
	if c == nil {
		c = DefaultConf
	}

	gwc := &GateWayConn{
		SendMsgs:   make(chan []byte),
		RecvMsgs:   make(chan []byte, 100),
		SendGWMsgs: make(chan []byte),
		RecvGWMsgs: make(chan []byte, 100),
		Configure:  c,
	}

	m = &GateWayManager{
		reportChan:    make(chan *Device, 100),
		Motions:       make(map[string]*Motion),
		Switchs:       make(map[string]*Switch),
		SensorHTs:     make(map[string]*SensorHT),
		Magnets:       make(map[string]*Magnet),
		Plugs:         make(map[string]*Plug),
		DiscoveryTime: time.Now().Unix(),
		GateWayConn:   gwc,
	}

	//connection
	err = gwc.initMultiCastConn(m.reportChan)
	if err != nil {
		return
	}

	//find the gateway
	m.whois()

	//discovery all device
	gw_ip := m.GateWayDevice.IP
	err = gwc.initGateWayConn(gw_ip)
	if err != nil {
		return
	}
	err = m.discovery()

	//report or heartbeat message
	go func() {
		for {
			dev := <-m.reportChan
			m.putDevice(dev)
		}
	}()

	return
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
		LOGGER.Warn("DEVICESYNC:: unknown model is %s", dev.Model)
	}

	return
}

func (m *GateWayManager) whois() {
	//read msg
	iamResp := &IamResp{}
	m.GateWayConn.communicate(NewWhoisRequest(), iamResp)

	//gateway infomation
	gwd := NewGateWayDevice(iamResp.Device)
	gwd.IP = iamResp.IP
	gwd.Port = iamResp.Port

	m.GateWayDevice = gwd
}

func (m *GateWayManager) discovery() (err error) {
	conn := m.GateWayConn

	//get devlist response
	LOGGER.Info("start to discover the device...")
	devListResp := &DeviceListResp{}
	if !conn.communicate(NewDevListRequest(), devListResp) {
		return errors.New("show device list error")
	}
	m.GateWayDevice.setToken(devListResp.Token)

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
