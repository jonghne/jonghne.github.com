package ws

import "time"

type client struct {
	id int64
	addr string
	send chan []byte
}
type Manager struct {
	links map[int64]client
	prot *Protocol
}

type EvHandle interface {
	OnUpdate(data []byte)
}

func CreateManager(prot *Protocol) *Manager {
	links := make(map[int64]client)
	manager := &Manager{links, prot}
	prot.Subscribe(manager)
	return manager
}

func (mgr *Manager) HandleMsg(id int64, msg []byte) {
	mgr.prot.Business()
}

func (mgr *Manager) Subscribe(send chan []byte) int64 {
	id := time.Now().UnixNano()
	user := client{id, "", send}
	mgr.links[id] = user
	return id
}


func (mgr *Manager) Unsubscribe(id int64) {
	delete(mgr.links, id)
}

func (mgr *Manager) OnUpdate(data []byte) {
	links := mgr.links
	for _, client := range links {
		go func() {
			client.send <- data
			client.send <- []byte("pong")
		}()
	}
}