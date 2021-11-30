package controllers

import (
	"conda/utils"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Add this lines
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var disconnect = make(chan bool)

func serve(connect *websocket.Conn) error {
	if err := connect.WriteMessage(websocket.TextMessage, []byte("hi.")); err != nil {
		log.Println(err)
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	defer func(connect *websocket.Conn) {
		ticker.Stop()
		if err := connect.Close(); err != nil {
			return
		}
	}(connect)
	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			if err := connect.WriteMessage(websocket.TextMessage, []byte(time.Now().String())); err != nil {
				log.Println(err)
				return err
			}
		case <-disconnect:
			ticker.Stop()
			if err := connect.Close(); err != nil {
				log.Println(err)
				return err
			}
			return nil
		}
	}
	return nil
}

func Connect(w http.ResponseWriter, r *http.Request) {
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := serve(connect); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
}

func Disconnect(w http.ResponseWriter, _ *http.Request) {
	select {
	case disconnect <- true:
		log.Println("Disconnect websocket")
	case <-time.After(1 * time.Second):
		log.Println("Timeout websocket")
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}
