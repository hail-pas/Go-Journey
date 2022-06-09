package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{receive: make(chan []byte, 256), ws: ws, info: &Data{}}
	hub.connectionsA <- c
	go c.writer()
	c.reader()
	defer func() {
		c.info.Type = "logout"
		userList = del(userList, c.info.User)
		c.info.UserList = userList
		c.info.Content = c.info.User
		dataB, _ := json.Marshal(c.info)
		hub.message <- dataB
		hub.connectionsA <- c
	}()
}

func Run() {
	router := mux.NewRouter()
	go hub.run()
	router.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}
