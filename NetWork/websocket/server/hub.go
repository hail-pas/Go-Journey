package server

import "encoding/json"

type hubStruct struct {
	connectionsAlive map[*connection]bool
	message          chan []byte
	connectionsA     chan *connection
	connectionsB     chan *connection
}

var hub = hubStruct{
	connectionsAlive: make(map[*connection]bool),
	connectionsB:     make(chan *connection),
	message:          make(chan []byte),
	connectionsA:     make(chan *connection),
}

func (hub *hubStruct) run() {
	for {
		select {
		case c := <-hub.connectionsA:
			hub.connectionsAlive[c] = true
			c.info.Ip = c.ws.RemoteAddr().String()
			c.info.Type = "handshake"
			c.info.UserList = userList
			data, _ := json.Marshal(c.info)
			c.receive <- data
		case c := <-hub.connectionsB:
			if _, ok := hub.connectionsAlive[c]; ok {
				delete(hub.connectionsAlive, c)
				close(c.receive)
			}
		case data := <-hub.message:
			for c := range hub.connectionsAlive {
				select {
				case c.receive <- data:
				default:
					delete(hub.connectionsAlive, c)
					close(c.receive)
				}
			}
		}
	}
}
