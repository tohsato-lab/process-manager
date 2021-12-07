package modules

import (
	"backend/repository"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/url"
)

var connections = make(map[string]*websocket.Conn)

func readPump(ip string, db *sqlx.DB) {
	for {
		_, message, err := connections[ip].ReadMessage()
		if err != nil {
			log.Println("read:", err)
			log.Println("disconnect backend")
			_ = repository.UpdateCalcServerStatus(db, ip, "stop")
			break
		}
		log.Printf("recv: %s", message)
	}
}

func Connection(ip string, port string, db *sqlx.DB) error {
	u := url.URL{Scheme: "ws", Host: ip + ":" + port, Path: "/connect"}
	log.Printf("connecting to %s", u.String())
	connect, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	connections[ip] = connect
	go readPump(ip, db)
	return nil
}

func Disconnection(ip string) error {
	if connections[ip] == nil {
		return nil
	}
	if err := connections[ip].WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		return err
	}
	if err := connections[ip].Close(); err != nil {
		return err
	}
	delete(connections, ip)
	return nil
}
