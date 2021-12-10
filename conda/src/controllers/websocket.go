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

func Connect(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	client := &modules.Client{DB: db, Conn: conn, Pipe: make(chan map[string]string, 256)}
	modules.Clients[client] = true
	go client.ReadPump()
	go client.WritePump()
}
