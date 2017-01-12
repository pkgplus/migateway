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

	CMD_IAM         = `iam`
	CMD_DEVLIST_ACK = `get_id_list_ack`
	CMD_HEARTBEAT   = `heartbeat`

	CMD_READ     = `read`
	CMD_READ_ACK = `read_ack`

	MSG_WHOIS   = `{"cmd":"whois"}`
	MSG_DEVLIST = `{"cmd":"get_id_list"}`
)

var (
	LOGGER = logs.GetBlogger()
)

type GateWayManager struct {
	GateWay          *GateWay
	Devices          map[string][]*Device
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

	err = initRead(m.RecvMsgs)
	if err != nil {
		return
	}
	err = initWriteMulti(m.SendMultiMsgs)
	if err != nil {
		return
	}

	//send whois
	m.SendMultiMsgs <- MSG_WHOIS

	//read msg
	iam_data := &GateWayDiscResp{}
	m.GetResp(CMD_IAM, iam_data)

	//gateway info
	m.GateWay = &GateWay{
		Sid:   iam_data.Sid,
		IP:    iam_data.IP,
		Model: iam_data.Model,
	}

	//init write to gateway
	err = initWriteGateWay(m.GateWay.IP, m.SendGWMsgs)
	if err != nil {
		return
	}

	//send devlist request
	m.SendGWMsgs <- MSG_DEVLIST

	//get devlist response
	var sids []string
	devlist_data := &GateWayDevListResp{}
	m.GetResp(CMD_DEVLIST_ACK, devlist_data)

	//get every device status
	for _, sid := range devlist_data.Data {
		m.SendGWMsgs <- NewGateWayReadRequest(sid)
		/////////////////////////m.GetResp(CMD_READ_ACK, resp)
	}

	return
}

func (m *GateWayManager) GetResp(cmd string, resp GateWayResp) {
	for {
		msg := <-m.RecvMsgs
		err = json.Unmarshal(msg, resp)
		if err != nil {
			return
		}

		if resp.GetCmd() == cmd {
			return
		} else {
			LOGGER.Warn("ingore the msg: %s", string(msg))
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
	UDP_ADDR = &net.UDPAddr{
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
