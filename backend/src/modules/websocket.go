package modules

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"time"
)

var connections map[string]*websocket.Conn

type Command struct {
	ProcessID string
	Command   string
}

func readPump(ip string) {
	for {
		_, message, err := connections[ip].ReadMessage()
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
	if connections == nil {
		connections = map[string]*websocket.Conn{}
	}
	connections[ip] = connect
	go readPump(ip)
	return nil
}

func Disconnection(ip string) error {
	if connections[ip] == nil {
		return nil
	}
	if err := connections[ip].WriteMessage(websocket.TextMessage, []byte("hi~~~~~")); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	if err := connections[ip].WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		return err
	}
	if err := connections[ip].Close(); err != nil {
		return err
	}
	delete(connections, ip)
	return nil
}
