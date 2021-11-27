package modules

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

var connections []*websocket.Conn

func sync(connect *websocket.Conn) {
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

func Connection(ip string, port string) error {
	u := url.URL{Scheme: "ws", Host: ip + ":" + port, Path: "/connect"}
	log.Printf("connecting to %s", u.String())
	connect, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	connections = append(connections, connect)
	go sync(connect)
	return nil
}
