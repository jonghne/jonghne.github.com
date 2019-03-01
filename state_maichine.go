package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"os"
	"net"
	"unsafe"
)

const (
	SERVER_IP       = "127.0.0.1"
	VOTE_NUM = 2
)

const (
	EV_TIMEOUT	= iota
	EV_CANDIDATE
	EV_LEADER
	EV_VOTE
)

type EVENT struct {
	ev int
	term int
	leader int
	value int
}

var SERVER_PORT = []int{10001,10002,10003,10004,10005,10006}
var adrList map[string]PEER

type PEER struct {
	addr *net.UDPAddr
	//conn net.Conn
}

type stateHandler interface {
	doAction(ctx *context, event interface{})
	info() string
}

type stateA struct {
	tm string
}

type stateB struct {
	tm string
}

type context struct {
	st stateHandler
	selfPort int
	conn *net.UDPConn
	msg chan EVENT
}

func (ctx *context) setState(st stateHandler) {
	ctx.st = st
}

func (ctx *context) getState() string {
	return ctx.st.info()
}

func (s *stateA) doAction(ctx *context, event interface{}) {
	fmt.Println("AAAAAAAAA", event)
	if v, ok := event.(uint32); ok {
		if v>5 {
			fmt.Println("change B")
			newState := &stateB{"stateB coming"}
			ctx.setState(newState)
		} else {
			fmt.Println("still in stateA")
		}
	}
}

func (s *stateA) info() string {
	return fmt.Sprint(s.tm, " working")
}

func (s *stateB) doAction(ctx *context, event interface{}) {
	fmt.Println("BBBBBBBBB", event)
	if v, ok := event.(uint32); ok {
		if v<6 {
			fmt.Println("change A")
			newState := &stateA{"stateA coming"}
			ctx.setState(newState)
		} else {
			fmt.Println("still in stateB")
		}
	}
}

func (s *stateB) info() string {
	return fmt.Sprint(s.tm, " working")
}

type sliceMock struct {
	addr uintptr
	len int
	cap int
}

func encode(msg EVENT) []byte {
	ep := unsafe.Pointer(&msg)
	tmp := &sliceMock{uintptr(ep),int(unsafe.Sizeof(msg)), int(unsafe.Sizeof(msg))}
	data := *(*[]byte)(unsafe.Pointer(tmp))
	return data
}

func decode(d []byte) EVENT {
	ep := *(**EVENT)(unsafe.Pointer(&d))
	return *ep
}

type STATE_LEADER struct {
	term int
}

func (s *STATE_LEADER) doAction(ctx *context, event interface{}) {
	if v, ok := event.(EVENT); ok {
		if v.ev == EV_CANDIDATE {
			if v.term > s.term {
				//chang to follower, vote
				newState := &STATE_FOLLOWER{v.term, v.leader}
				voteMsg := EVENT{EV_VOTE, v.term, v.leader, 0}
				data := encode(voteMsg)
				ctx.broadcastData(data)
				ctx.setState(newState)
			}
		} else if v.ev == EV_TIMEOUT {
			//chang to follower
			newState := &STATE_FOLLOWER{s.term, 0}
			ctx.setState(newState)
		}
	}
}

func (s *STATE_LEADER) info() string {
	return "i am leader"
}

type STATE_CANDIDATE struct {
	term int
	vote int
}

func (s *STATE_CANDIDATE) doAction(ctx *context, event interface{}) {
	if v, ok := event.(EVENT); ok {
		if v.ev == EV_VOTE {
			if v.leader == ctx.selfPort {
				s.vote++
				if s.vote >= VOTE_NUM {
					newState := &STATE_LEADER{s.term}
					ctx.setState(newState)
				}
			} else if v.term > s.term {
				//chang to follower, vote
				newState := &STATE_FOLLOWER{v.term, v.leader}
				voteMsg := EVENT{EV_VOTE, v.term, v.leader, 0}
				data := encode(voteMsg)
				ctx.broadcastData(data)
				ctx.setState(newState)
			}
		} else if v.ev == EV_CANDIDATE {
			if v.term > s.term {
				//chang to follower, vote
				newState := &STATE_FOLLOWER{v.term, v.leader}
				voteMsg := EVENT{EV_VOTE, v.term, v.leader, 0}
				data := encode(voteMsg)
				ctx.broadcastData(data)
				ctx.setState(newState)
			}
		} else if v.ev == EV_TIMEOUT {
			//elect
			s.term++
			s.vote = 0
			electMsg := EVENT{EV_CANDIDATE, s.term, ctx.selfPort, 0}
			data := encode(electMsg)
			ctx.broadcastData(data)
		}
	}
}

func (s *STATE_CANDIDATE) info() string {
	return "i am candidate"
}

type STATE_FOLLOWER struct {
	term int
	leader int
}

func (s *STATE_FOLLOWER) doAction(ctx *context, event interface{}) {
	if v, ok := event.(EVENT); ok {
		if v.ev == EV_CANDIDATE {
			if v.term > s.term {
				//vote
				s.term = v.term
				s.leader = v.leader
				voteMsg := EVENT{EV_VOTE, v.term, v.leader, 0}
				data := encode(voteMsg)
				ctx.broadcastData(data)
			}
		} else if v.ev == EV_TIMEOUT {
			//change to candidate, elect
			newState := &STATE_CANDIDATE{s.term+1, 0}
			electMsg := EVENT{EV_CANDIDATE, s.term+1, ctx.selfPort, 0}
			data := encode(electMsg)
			ctx.broadcastData(data)
			ctx.setState(newState)
		}
	}
}

func (s *STATE_FOLLOWER) info() string {
	return "i am follower"
}

func (ctx *context) start() {
	timer := time.NewTimer(10*time.Second)
	//receive peers' messages
	go ctx.recvMsg()
	for {
		select {
		case msg := <- ctx.msg:
			ctx.st.doAction(ctx, msg)
		case <- timer.C:
			timeoutEv := EVENT{}
			timeoutEv.ev = EV_TIMEOUT
			ctx.msg <- timeoutEv
		}
		timer.Reset(10*time.Second)
	}
}

func initialize() (ctx *context) {
	ctx = &context{}
	initState := &STATE_FOLLOWER{0, 0}
	ctx.st = initState
	ctx.msg = make(chan EVENT, 6)
	return ctx
}

func (ctx *context) sendData(to *net.UDPAddr, data []byte) (int, error) {
	n, err := ctx.conn.WriteToUDP(data, to)
	return n, err
}

func (ctx *context) broadcastData(data []byte) {
	for _, peer := range adrList {
		ctx.sendData(peer.addr, data)
	}
}

func (ctx *context) recvMsg() {
	var buf [64]byte
	conn := ctx.conn

	for {
		n, rAddr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}
		if _, ok := adrList[rAddr.String()]; !ok {
			adrList[rAddr.String()] = PEER{rAddr}
		}

		event := decode(buf[:n])
		fmt.Println("from ", rAddr, "recv ", event)

		ctx.msg <- event
	}
}

func (ctx *context) scanPeer() {
	for i:=0; i<len(SERVER_PORT); i++ {
		if SERVER_PORT[i] != ctx.selfPort {
			address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT[i])
			conn, err := net.Dial("udp", address)
			if err == nil {
				target, e := net.ResolveUDPAddr("udp", address)
				if e != nil {
					fmt.Println(e)
					os.Exit(1)
				}
				adrList[conn.RemoteAddr().String()] = PEER{target}
				fmt.Println(address, "already online")
				conn.Close()
			}
		}
	}
}

func (ctx *context) createListener() {
	rand.Seed(time.Now().UnixNano())
	no := rand.Uint32() % 6
	ctx.selfPort = SERVER_PORT[no]
	address := SERVER_IP + ":" + strconv.Itoa(ctx.selfPort)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx.scanPeer()
}

func init() {
	adrList = make(map[string]PEER)
}

func main() {
	//init stateA
	//s := &stateA{"stateA init"}
	//ctx := &context{}
	//for {
	//	rand.Seed(time.Now().UnixNano())
	//	ev := rand.Uint32()%10
	//	ctx.st.doAction(ctx, ev)
	//	fmt.Println(ctx.getState())
	//	time.Sleep(time.Second)
	//}

	ctx := initialize()
	ctx.createListener()
	ctx.start()
}