package modules

import (
	"backend/repository"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"net/url"
)

var connections = make(map[string]*websocket.Conn)

func readPump(ip string, db *sqlx.DB) error {
	for {
		_, message, err := connections[ip].ReadMessage()
		if err != nil {
			log.Println("read:", err)
			log.Println("disconnect backend")
			if err := repository.UpdateCalcServerStatus(db, ip, "stop"); err != nil {
				return err
			}
			break
		}
		log.Printf("recv: %s", message)
		var contents map[string]string
		if err := json.Unmarshal(message, &contents); err != nil {
			return err
		}
		switch contents["status"] {
		case "ready":
		case "running":
			if err := repository.UpdateProcessStatus(db, contents["ID"], "running"); err != nil {
				return err
			}
			if err := repository.SetStartDate(db, contents["ID"]); err != nil {
				return err
			}
		default:
			if err := repository.UpdateProcessStatus(db, contents["ID"], "syncing"); err != nil {
				return err
			}
			if err := repository.SetCompleteDate(db, contents["ID"]); err != nil {
				return err
			}
			go func() {
				log.Println("sync")
				_, err := Rsync(db, contents["ID"])
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("sync done")
				if err := repository.UpdateProcessStatus(db, contents["ID"], contents["status"]); err != nil {
					return
				}
				if err := UpdateProcess(db); err != nil {
					return
				}
			}()
		}
		if err := UpdateProcess(db); err != nil {
			return err
		}
	}
	return nil
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
	go func() {
		if err := readPump(ip, db); err != nil {
			log.Println(err)
			return
		}
	}()
	processIDs, err := repository.NeedSyncProcesses(db, ip)
	log.Println(processIDs)
	if err != nil {
		return err
	}
	var commands []map[string]string
	for _, processID := range processIDs {
		commands = append(commands, map[string]string{"ID": processID, "status": "sync"})
	}
	if err := connections[ip].WriteJSON(commands); err != nil {
		return err
	}
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
