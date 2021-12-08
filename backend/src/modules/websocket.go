package modules

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/url"

	"backend/repository"
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
		var contents map[string]string
		if err := json.Unmarshal(message, &contents); err != nil {
			log.Println(err)
			return
		}
		if err := repository.UpdateProcessStatus(db, contents["ID"], contents["status"]); err != nil {
			log.Println(err)
			return
		}
		switch contents["status"] {
		case "ready":
		case "running":
			if err := repository.SetStartDate(db, contents["ID"]); err != nil {
				log.Println(err)
				return
			}
		default:
			if err := repository.SetCompleteDate(db, contents["ID"]); err != nil {
				log.Println(err)
				return
			}
		}
		go UpdateProcess(db)
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
