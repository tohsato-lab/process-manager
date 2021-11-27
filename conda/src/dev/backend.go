package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

func register(w http.ResponseWriter, r *http.Request) {
	go func() {
		u := url.URL{Scheme: "ws", Host: ":8100", Path: "/connect", RawQuery: "value=test"}
		log.Printf("connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		defer func(c *websocket.Conn) {
			if c.Close() != nil {
				return
			}
		}(c)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				log.Println("disconnect")
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}

func main() {
	http.HandleFunc("/register", register)
	log.Fatal(http.ListenAndServe(":8101", nil))
}
