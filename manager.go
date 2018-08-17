package migateway

import (
	"errors"
	"time"

	"github.com/bingbaba/util/logs"
)

var (
	LOGGER ILogger = logs.GetBlogger()
	EOF    byte    = 0
)

type AqaraManager struct {
	GateWay           *GateWay
	Motions           map[string]*Motion
	Switchs           map[string]*Switch
	DualWiredSwitches map[string]*DualWiredWallSwitch
	SensorHTs         map[string]*SensorHT
	Magnets           map[string]*Magnet
	Plugs             map[string]*Plug
	ReportAllMessages bool

	StateMessages chan interface{}

	DiscoveryTime    int64
	FreshDevListTime int64
}

func NewAqaraManager() (m *AqaraManager) {
	//AqaraManager
	m = &AqaraManager{
		Motions:           make(map[string]*Motion),
		Switchs:           make(map[string]*Switch),
		DualWiredSwitches: make(map[string]*DualWiredWallSwitch),
		SensorHTs:         make(map[string]*SensorHT),
		Magnets:           make(map[string]*Magnet),
		Plugs:             make(map[string]*Plug),
		StateMessages:     make(chan interface{}, 1024),
		DiscoveryTime:     time.Now().Unix(),
	}

	return
}

func (m *AqaraManager) Start(c *Configure) (err error) {
	if c == nil {
		c = DefaultConf
	}

	m.ReportAllMessages = c.ReportAllMessages

	// Connection
	conn := NewConn(c)
	err = conn.initMultiCast()
	if err != nil {
		return
	}

	// Find gateway
	m.whois(conn)

	// Show device list
	gw_ip := m.GateWay.IP
	err = conn.initGateWay(gw_ip)
	if err != nil {
		return
	}

	err = m.discovery()

	// Report or heartbeat message
	go func() {
		for msg := range conn.devMsgs {
			m.putDevice(msg)
		}
	}()

	return
}

func (m *AqaraManager) Stop() {
	if nil != m.GateWay.GatewayConnection {
		m.GateWay.GatewayConnection.stop()
		m.GateWay.GatewayConnection.closeRead <- true
		m.GateWay.GatewayConnection.closeWrite <- true
	}

	close(m.GateWay.GatewayConnection.devMsgs)
	close(m.StateMessages)
}

func (m *AqaraManager) putDevice(dev *Device) (added bool) {
	if nil == dev {
		return false
	}

	LOGGER.Info("DEVICESYNC:: %s(%s): %s", dev.Model, dev.Sid, dev.Data)
	gateway := m.GateWay

	var saveDev *Device
	added = true
	switch dev.Model {
	case MODEL_GATEWAY:
		gateway.Set(dev)
		saveDev = gateway.Device
	case MODEL_MOTION:
		d, found := m.Motions[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.Motions[dev.Sid] = NewMotion(dev)
		}
		saveDev = m.Motions[dev.Sid].Device
	case MODEL_SWITCH:
		d, found := m.Switchs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.Switchs[dev.Sid] = NewSwitch(dev)
		}
		saveDev = m.Switchs[dev.Sid].Device
	case MODEL_DUALWIREDSWITCH:
		d, found := m.DualWiredSwitches[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.DualWiredSwitches[dev.Sid] = NewDualWiredSwitch(dev)
		}
		saveDev = m.DualWiredSwitches[dev.Sid].Device
	case MODEL_SENSORHT:
		d, found := m.SensorHTs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.SensorHTs[dev.Sid] = NewSensorHt(dev)
		}
		saveDev = m.SensorHTs[dev.Sid].Device
	case MODEL_MAGNET:
		d, found := m.Magnets[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.Magnets[dev.Sid] = NewMagnet(dev)
		}
		saveDev = m.Magnets[dev.Sid].Device
	case MODEL_PLUG:
		d, found := m.Plugs[dev.Sid]
		if found {
			d.Set(dev)
		} else {
			dev.Aqara = m
			dev.GatewayConnection = gateway.GatewayConnection
			m.Plugs[dev.Sid] = NewPlug(dev)
		}
		saveDev = m.Plugs[dev.Sid].Device
	default:
		added = false
		LOGGER.Warn("DEVICESYNC:: unknown model is %s", dev.Model)
	}

	LOGGER.Debug("save to report chan...")
	if saveDev != nil {

	}
	LOGGER.Debug("save to report chan over!")

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
	dev.Aqara = m
	dev.GatewayConnection = conn

	m.GateWay = dev
}

func (m *AqaraManager) discovery() (err error) {
	gateway := m.GateWay
	conn := gateway.GatewayConnection

	//get devlist response
	LOGGER.Info("start to discover the device...")
	devListResp := &DeviceListResp{}
	if !conn.communicate(NewDevListRequest(), devListResp) {
		return errors.New("show device list error")
	}
	//gateway.setToken(devListResp.Token)
	gateway.GatewayConnection.token = devListResp.Token

	//every device
	for index, sid := range devListResp.getSidArray() {
		LOGGER.Info("DISCOVERY[%d] of device %s", index, sid)
		dev := conn.waitDevice(sid)
		dev.Aqara = m
		dev.GatewayConnection = conn
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
	m.GateWay.GatewayConnection.SetAESKey(key)
}
