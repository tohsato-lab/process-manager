package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	"backend/modules"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Connect(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println(err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	client := &modules.Client{DB: db, Conn: conn, Send: make(chan []byte, 256)}
	modules.SocketCore.Register <- client
	go client.WritePump()
}
