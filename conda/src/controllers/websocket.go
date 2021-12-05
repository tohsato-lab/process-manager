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

func readPump(connect *websocket.Conn) {
	defer func(c *websocket.Conn) {
		if c.Close() != nil {
			return
		}
	}(connect)
	for {
		_, message, err := connect.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			log.Println("disconnect")
			return
		}
		log.Printf("recv: %s", message)
	}
}

func writePump(connect *websocket.Conn) {
	defer func(c *websocket.Conn) {
		if c.Close() != nil {
			return
		}
	}(connect)
	for {
		select {
		case cmd := <-Command:
			if err := connect.WriteJSON(cmd); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func serve(connect *websocket.Conn) error {
	if err := connect.WriteMessage(websocket.TextMessage, []byte("hi.")); err != nil {
		log.Println(err)
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	defer func(connect *websocket.Conn) {
		log.Println("disconnect1")
		ticker.Stop()
		if err := connect.Close(); err != nil {
			return
		}
	}(connect)
	for {
		select {
		case <-ticker.C:
			/*
				if err := connect.WriteMessage(websocket.TextMessage, []byte(time.Now().String())); err != nil {
						log.Println(err)
						return err
				}
			*/
		case <-disconnect:
			log.Println("disconnect2")
			return nil
		}
	}
}

func Connect(w http.ResponseWriter, r *http.Request) {
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	go serve(connect)
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
