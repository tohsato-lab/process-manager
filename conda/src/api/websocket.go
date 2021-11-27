package api

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func Connect(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	if c.WriteMessage(websocket.TextMessage, []byte("hi.")) != nil {
		log.Print("upgrade:", err)
		return
	}
	ticker := time.NewTicker(1 * time.Second)
	defer func() {
		ticker.Stop()
		if c.Close() != nil {
			return
		}
	}()
	for {
		select {
		case <-ticker.C:
			if c.WriteMessage(websocket.TextMessage, []byte(time.Now().String())) != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
