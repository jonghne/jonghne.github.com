package main

import (
	"ws"
)
func main() {
	prot := ws.CreateProtocol()
	manager := ws.CreateManager(prot)
	ws.StartCommunicate(manager)
}
