package migateway

import (
	"encoding/json"
	"errors"
	"net"
	"time"
)

const (
	MULTICAST_ADDR = "224.0.0.50"
	MULTICAST_PORT = 4321
	SERVER_PORT    = 9898
)

type GateWayConn struct {
	conn       *net.UDPConn
	closeRead  chan bool
	closeWrite chan bool
	SendMsgs   chan []byte
	RecvMsgs   chan []byte
	SendGWMsgs chan []byte
	RecvGWMsgs chan []byte
	devMsgs    chan *Device
	*Configure

	token string
}

func NewConn(c *Configure) *GateWayConn {
	return &GateWayConn{
		closeRead:  make(chan bool),
		closeWrite: make(chan bool),
		SendMsgs:   make(chan []byte),
		RecvMsgs:   make(chan []byte, 100),
		SendGWMsgs: make(chan []byte),
		RecvGWMsgs: make(chan []byte, 100),
		devMsgs:    make(chan *Device, 100),
		Configure:  c,
	}
}

func (gwc *GateWayConn) Control(dev *Device, data map[string]interface{}) error {
	// wait for multicast heartbeat message from device
	if !dev.waitToken() {
		return errors.New("ERROR: heartbeat timeout")
	}

	req, err := newWriteRequest(dev, gwc.Configure.AESKey, data)
	if err != nil {
		return err
	}
	gwc.Send(req)
	return nil
}

func (gwc *GateWayConn) Send(req *Request) bool {
	if req == nil {
		return false
	}

	if req.Cmd == CMD_WHOIS {
		gwc.multicast(req)
	} else {
		gwc.sendGW(req)
	}
	return true
}

func (gwc *GateWayConn) waitDevice(sid string) *Device {
	req := NewReadRequest(sid)
	resp := &DeviceBaseResp{}
	for {
		if gwc.communicate(req, resp) {
			if resp.Sid == sid {
				break
			} else {
				LOGGER.Info("get a unknown device %s, model = %s, sid = %s", CMD_READ_ACK, resp.Model, resp.Sid)
			}
		}
	}
	return resp.Device
}

func (gwc *GateWayConn) initMultiCast() error {
	//listen
	udpaddr := &net.UDPAddr{IP: net.ParseIP(MULTICAST_ADDR), Port: SERVER_PORT}
	con, err := net.ListenMulticastUDP("udp4", nil, udpaddr)
	if err != nil {
		return err
	}
	LOGGER.Info("listening %d ...", SERVER_PORT)

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
				//  LOGGER.Warn("MULTICAST:: parse invalid msg: %s, error:%v", buf[0:size], errTmp)
				// }
				if resp.Cmd == CMD_REPORT {
					gwc.devMsgs <- resp.Device
				} else if resp.Cmd == CMD_HEARTBEAT {
					resp.freshHeartTime()
					gwc.devMsgs <- resp.Device
				} else {
					msg := make([]byte, size)
					copy(msg, buf[0:size])
					gwc.RecvMsgs <- msg
				}

			}
		}
	}()

	//write
	udpaddrMulticast := &net.UDPAddr{
		IP:   net.ParseIP(MULTICAST_ADDR),
		Port: MULTICAST_PORT,
	}
	go func() {
		for {
			msg := <-gwc.SendMsgs
			wsize, err3 := con.WriteToUDP(msg, udpaddrMulticast)
			if err3 != nil {
				panic(err3)
			}
			LOGGER.Info("MULTICAST:: send msg: %s, %d bytes!", msg, wsize)
		}
	}()

	return nil
}

/*func (gwc *GateWayConn) forwardReport(dev *Device) bool {
	var forwarded bool
	select {
	case gwc.ReportMsgs <- dev:
		forwarded = true
	case <-time.After(time.Duration(gwc.ReportForwardTimeout) * time.Second):
		LOGGER.Error("forward report message timeout:: %s", string(toBytes(dev)))
	}
	return forwarded
}
*/

func (gwc *GateWayConn) multicast(req *Request) {
	gwc.SendMsgs <- toBytes(req)
}

func (gwc *GateWayConn) sendGW(req *Request) {
	gwc.SendGWMsgs <- toBytes(req)
}

func (gwc *GateWayConn) communicate(req *Request, resp Response) bool {
	expectCmd := req.expectCmd()
	if expectCmd == "" {
		LOGGER.Warn("unknown request: %s", string(toBytes(req)))
		return false
	}
	//send message
	gwc.Send(req)

	retry := 0
	maxRetry, timeout := gwc.Configure.getRetryAndTimeout(req)
	chanName := req.getChanName()

	LOGGER.Info("%s:: wait \"%s\" response...", chanName, expectCmd)
	for {
		select {
		case msg := <-gwc.getChan(req.Cmd):
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
		case <-time.After(timeout):
			retry++
			if retry > maxRetry {
				LOGGER.Error("%s:: recv msg TIMEOUT", chanName)
				return false
			} else {
				LOGGER.Error("%s:: send msg retry %d ...", chanName, retry)
				gwc.Send(req)
			}
		}
	}
	return false
}

func (gwc *GateWayConn) initGateWay(ip string) (err error) {

	//do we first need to close an existing connection ?
	if gwc.conn != nil {
		// Signal the read and write go-routine to terminate
		gwc.closeRead <- true
		gwc.closeWrite <- true
		gwc.conn.Close()
	}

	UDPAddr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: SERVER_PORT,
	}

	//open new conn
	gwc.conn, err = net.DialUDP("udp4", nil, UDPAddr)
	if err != nil {
		return
	}

	//write
	go func() {
		for {
			select {
			case <-gwc.closeWrite:
				return
			case msg := <-gwc.SendGWMsgs:
				LOGGER.Info("GATEWAY:: send msg: %s", msg)
				_, werr := gwc.conn.Write([]byte(msg))
				if werr != nil {
					LOGGER.Error("send error %v", werr)
				}
			}
		}
	}()

	//read
	go func() {
		buf := make([]byte, 2048)
		for {
			select {
			case <-gwc.closeRead:
				return
			default:
				size, _, err2 := gwc.conn.ReadFromUDP(buf)
				if err2 != nil {
					//panic(err2)
					LOGGER.Error("GATEWAY:: recv error: %v", err2)
				} else if size > 0 {
					LOGGER.Debug("GATEWAY:: recv msg: %s", string(buf[0:size]))
					msg := make([]byte, size)
					copy(msg, buf[0:size])
					gwc.RecvGWMsgs <- msg
				}
			}
		}
	}()

	return
}

func (gwc *GateWayConn) getChan(cmd string) chan []byte {
	if cmd == CMD_WHOIS {
		return gwc.RecvMsgs
	}
	return gwc.RecvGWMsgs
}
