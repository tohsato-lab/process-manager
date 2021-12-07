package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"net/http"

	"conda/modules"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Connect(w http.ResponseWriter, r *http.Request, hub *modules.Hub, db *sqlx.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	client := &modules.Client{Hub: hub, DB: db, Conn: conn, Pipe: make(chan string, 256)}
	client.Hub.Register <- client
	go client.ReadPump()
	go client.WritePump()
}
