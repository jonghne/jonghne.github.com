package ws

import (
	"fmt"
	"math/rand"
	"time"
)

type Protocol struct {
	ev EvHandle
}


func CreateProtocol() *Protocol {
	return &Protocol{}
}


func (prot *Protocol) Subscribe(ev EvHandle) {
	prot.ev = ev
}

func (prot *Protocol) Business() {
	rand.Seed(time.Now().UnixNano())
	num := rand.Int()
	prot.ev.OnUpdate([]byte(fmt.Sprintf("update %d", num)))
}