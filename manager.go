package migateway

import (
	"encoding/json"
	"github.com/bingbaba/util/logs"
	"net"
	"time"
)

const (
	MULTICAST_ADDR = "224.0.0.50"
	MULTICAST_PORT = 4321
	SERVER_PORT    = 9898

	MSG_WHOIS   = `{"cmd":"whois"}`
	MSG_DEVLIST = `{"cmd":"get_id_list"}`
)

var (
	LOGGER      = logs.GetBlogger()
	EOF    byte = 0
)

type GateWayManager struct {
	*GateWayDevice

	IP        string
	Port      string
	Motions   map[string]*Motion
	Switchs   map[string]*Switch
	SensorHTs map[string]*SensorHT

	DiscoveryTime    int64
	FreshDevListTime int64

	SendMsgs   chan string
	RecvMsgs   chan []byte
	SendGWMsgs chan string
	RecvGWMsgs chan []byte
}

func NewGateWayManager() (m *GateWayManager, err error) {
	m = &GateWayManager{
		Motions:       make(map[string]*Motion),
		Switchs:       make(map[string]*Switch),
		SensorHTs:     make(map[string]*SensorHT),
		DiscoveryTime: time.Now().Unix(),
		SendMsgs:      make(chan string),
		RecvMsgs:      make(chan []byte, 100),
		SendGWMsgs:    make(chan string),
		RecvGWMsgs:    make(chan []byte, 100),
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

func (m *GateWayManager) whois() {
	//send whois
	m.SendMsgs <- MSG_WHOIS

	//read msg to gateway
	iamResp := m.waitIam()
	m.GateWayDevice = &GateWayDevice{
		Device: iamResp.Device,
	}
	m.IP = iamResp.IP
	m.Port = iamResp.Port
}

func (m *GateWayManager) discovery() (err error) {
	LOGGER.Info("start to discover the device...")

	//send devlist request
	m.SendGWMsgs <- MSG_DEVLIST

	//get devlist response
	devListResp := m.waitDevList()
	for index, sid := range devListResp.getSidArray() {
		m.SendGWMsgs <- NewReadRequest(sid)

		dev := m.WaitDevice(sid)
		if m.addDevice(dev) {
			LOGGER.Warn("DISCOVERY[%d]: found the device %s: %v", index, dev.Model, dev.Data)
		} else {
			LOGGER.Warn("DISCOVERY[%d]: unknown model %s device: %v", index, dev.Model, dev)
		}
	}
	m.DiscoveryTime = time.Now().UnixNano()

	return
}

func (m *GateWayManager) addDevice(dev *Device) (added bool) {
	added = true

	switch dev.Model {
	case MODEL_MOTION:
		m.Motions[dev.Sid] = NewMotion(dev)
	case MODEL_SWITCH:
		m.Switchs[dev.Sid] = NewSwitch(dev)
	case MODEL_SENSORHT:
		m.SensorHTs[dev.Sid] = NewSensorHt(dev)
	default:
		added = false
	}

	return
}

func (m *GateWayManager) waitIam() *IamResp {
	resp := &IamResp{}
	m.waitResp(CMD_IAM, resp)
	return resp
}

func (m *GateWayManager) waitDevList() *DeviceListResp {
	resp := &DeviceListResp{}
	m.waitResp(CMD_DEVLIST_ACK, resp)
	return resp
}

func (m *GateWayManager) WaitDevice(sid string) *Device {
	resp := &DeviceBaseResp{}
	for {
		m.waitResp(CMD_READ_ACK, resp)
		if resp.Sid == sid {
			break
		} else {
			LOGGER.Info("get a unknown device %s, model = %s, sid = %s", CMD_READ_ACK, resp.Model, resp.Sid)
		}
	}
	return resp.Device
}

func (m *GateWayManager) waitResp(cmd string, resp Response) {
	LOGGER.Info("wait \"%s\" response...", cmd)

	var chanName string
	var msgChan chan []byte
	if cmd != CMD_IAM && cmd != CMD_REPORT {
		msgChan = m.RecvGWMsgs
		chanName = "GATEWAY"
	} else {
		msgChan = m.RecvMsgs
		chanName = "MULTICAST"
	}

	for {
		msg := <-msgChan
		err := json.Unmarshal(msg, resp)
		if err != nil {
			LOGGER.Error("%s:: parse %s error: %v", chanName, string(msg), err)
			continue
		} else if resp.GetCmd() != cmd {
			LOGGER.Warn("%s:: wait %s, ingore the msg: %s", chanName, cmd, string(msg))
			continue
		} else {
			LOGGER.Info("%s:: recv msg: %s", chanName, string(msg))
			break
		}
	}
	return
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
				g.RecvMsgs <- buf[0:size]
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
			wsize, err3 := con.WriteToUDP([]byte(MSG_WHOIS), MULTI_UDP_ADDR)
			if err3 != nil {
				panic(err3)
			}
			LOGGER.Info("MULTICAST:: send msg: %s, %d bytes!", msg, wsize)
		}
	}()

	return nil
}
