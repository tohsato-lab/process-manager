package api

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func serve(connect *websocket.Conn) {
	if connect.WriteMessage(websocket.TextMessage, []byte("hi.")) != nil {
		return
	}
	ticker := time.NewTicker(1 * time.Second)
	defer func(connect *websocket.Conn) {
		ticker.Stop()
		if connect.Close() != nil {
			return
		}
	}(connect)
	for i := 0; i < 5; i++ {
		select {
		case <-ticker.C:
			if connect.WriteMessage(websocket.TextMessage, []byte(time.Now().String())) != nil {
				return
			}
		}
	}
}

func Connect(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusVariantAlsoNegotiates)
		return
	}
	serve(connect)
	return
}
