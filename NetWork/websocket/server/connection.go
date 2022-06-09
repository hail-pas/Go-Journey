package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	ws      *websocket.Conn
	receive chan []byte
	info    *Data
}

var wu = &websocket.Upgrader{ReadBufferSize: 512, WriteBufferSize: 512, CheckOrigin: func(r *http.Request) bool {
	return true
}}

func (c *connection) writer() {
	for message := range c.receive {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
	err := c.ws.Close()
	if err != nil {
		return
	}
}

var userList []string

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			hub.connectionsA <- c
			break
		}
		json.Unmarshal(message, &c.info)
		switch c.info.Type {
		case "login":
			c.info.User = c.info.Content
			c.info.From = c.info.User
			userList = append(userList, c.info.User)
			c.info.UserList = userList
			dataB, _ := json.Marshal(c.info)
			hub.message <- dataB
		case "user":
			c.info.Type = "user"
			dataB, _ := json.Marshal(c.info)
			hub.message <- dataB
		case "logout":
			c.info.Type = "logout"
			userList = del(userList, c.info.User)
			dataB, _ := json.Marshal(c.info)
			hub.message <- dataB
			hub.connectionsA <- c
		default:
			fmt.Print("========default================")
		}
	}
}

func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var nSlice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			nSlice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(nSlice)
	return nSlice
}
