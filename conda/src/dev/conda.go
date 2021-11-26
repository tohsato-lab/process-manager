package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options

func connect(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.FormValue("value"))
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

func main() {
	http.HandleFunc("/connect", connect)
	log.Fatal(http.ListenAndServe(":8100", nil))
}
