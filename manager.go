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

	MSG_WHOIS   = `{"cmd":"whois"}`
	MSG_DEVLIST = `{"cmd":"get_id_list"}`
)

var (
	LOGGER = logs.GetBlogger()
)

type GateWayManager struct {
	*GateWayDevice
	IP               string
	Port             string
	Devices          map[string][]Device
	DiscoveryTime    int64
	FreshDevListTime int64
	SendMsgs         chan string
	RecvMsgs         chan string
	SendGWMsgs       chan string
}

func NewGateWayManager() (m *GateWayManager, err error) {
	m = &GateWayManager{
		Devices:       make(map[string][]*Device),
		DiscoveryTime: time.Now().Unix(),
		SendMultiMsgs: make(chan string),
		RecvMsgs:      make(chan []byte, 100),
		SendGWMsgs:    make(chan string),
	}

	//connection
	err = initRead(m.RecvMsgs)
	if err != nil {
		return
	}
	err = initWriteMulti(m.SendMultiMsgs)
	if err != nil {
		return
	}

	//find the gateway
	m.whois()

	//discovery all device
	err = m.discovery()

	return
}

func (m *GateWayManager) whois() {
	//send whois
	m.SendMultiMsgs <- MSG_WHOIS

	//read msg to gateway
	iamResp := m.WaitIam()
	m.GateWayDevice = &GateWayDevice{
		DeviceBaseInfo: iamResp.DeviceBaseInfo,
	}
	m.IP = iamResp.IP
	m.Port = iamResp.Port
}

func (m *GateWayManager) discovery() error {
	LOGGER.Info("start to discover the device...")

	//init write to gateway
	err := initWriteGateWay(m.IP, m.SendGWMsgs)
	if err != nil {
		return
	}

	//send devlist request
	m.SendGWMsgs <- MSG_DEVLIST

	//get devlist response
	devListResp := m.WaitDevList()
	for index, sid := range devListResp.Data {
		m.SendGWMsgs <- NewReadRequest(sid)
		dev := m.WaitDevice(nil)

		devs, found := m.Devices[dev.GetModel()]
		if !found {
			devs = make(map[string][]Device, 0)
		}
		devs = append(devs, dev)
		m.Devices[dev.GetModel()] = devs
		LOGGER.Info("DISCOVERY[%d] found the device %s: %v", index, dev.GetModel(), dev.GetData())
	}

	m.DiscoveryTime = time.Now().UnixNano()
	return
}

func (m *GateWayManager) WaitIam() *IamResp {
	resp := &IamResp{}
	waitResp(CMD_IAM, resp)
	return resp
}

func (m *GateWayManager) WaitDevList() *DeviceListResp {
	resp := &DeviceListResp{}
	waitResp(CMD_DEVLIST_ACK, resp)
	return resp
}

func (m *GateWayManager) WaitDevice(model string) Device {
	//model is empty
	if model == nil || model == "" {
		dev := &DeviceBaseInfo{}
		m.waitResp(dev)
		LOGGER.Info("get model %s", dev.GetModel())
		return WaitDevice(dev.GetModel())
	} else {
		var dev Device
		switch model {
		case DEV_MODEL_GATEWAY:
			dev = &GateWayDevice{}
		case MODEL_MOTION:
			dev = &MotionDevice{}
		default:
			LOGGER.Warn("unknown model %s", model)
			dev = &DeviceBaseInfo{}
		}
		m.waitResp(dev)
		dev.freshHeartTime()
		return dev
	}
}

func (m *GateWayManager) waitResp(cmd string, resp Response) {
	LOGGER.Info("wait \"%s\" response...", cmd)
	for {
		msg := <-m.RecvMsgs
		err = json.Unmarshal(msg, resp)
		if err != nil {
			LOGGER.Error("parse %s error: %v", string(msg), err)
			continue
		} else if resp.GetCmd() != cmd {
			LOGGER.Warn("wait %s, ingore the msg: %s", cmd, string(msg))
			continue
		} else {
			break
		}
	}
}

func initWriteMulti(msgs chan string) error {
	MULTI_UDP_ADDR = &net.UDPAddr{
		IP:   net.ParseIP(MULTICAST_ADDR),
		Port: MULTICAST_PORT,
	}
	mCon, err := net.DialUDP("udp", nil, MULTI_UDP_ADDR)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg := <-msgs
			LOGGER.Info("MULTICAST:: send msg: %s", msg)
			_, werr := mCon.Write([]byte(msg))
			if werr != nil {
				LOGGER.Error("send error %v", werr)
			}
		}
	}()
}

func initWriteGateWay(ip string, msgs chan string) error {
	UDP_ADDR := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: SERVER_PORT,
	}
	mCon, err := net.DialUDP("udp", nil, UDP_ADDR)
	if err != nil {
		return err
	}

	go func() {
		for {
			msg := <-msgs
			LOGGER.Info("GATEWAY:: send msg: %s", msg)
			_, werr := mCon.Write([]byte(msg))
			if werr != nil {
				LOGGER.Error("send error %v", werr)
			}
		}
	}()
}

func initRead(msgs chan string) error {
	udp_l := &net.UDPAddr{IP: net.ParseIP(MULTICAST_ADDR), Port: SERVER_PORT}
	con, err := net.ListenMulticastUDP(network, nil, udp_l)
	if err != nil {
		return err
	}
	LOGGER.Info("listennig %d ...", SERVER_PORT)

	go func() {
		defer con.Close()

		buf := make([]byte, 2048)
		for {
			_, _, err2 := con.ReadFromUDP(buf)
			if err2 != nil {
				panic(err2)
			} else {
				LOGGER.Debug("MULTICAST:: recv msg: %s", string(buf))
				msgs <- buf
			}
		}
	}()

	return nil
}
