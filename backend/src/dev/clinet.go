package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var clients = make(map[*websocket.Conn]bool)

func connection() **websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ":8100", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(c *websocket.Conn) {
		if c.Close() != nil {
			log.Println("disconnect")
			return
		}
	}(c)
	return &c
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for i := 0; i < 3; i++ {
		u := url.URL{Scheme: "ws", Host: ":8100", Path: "/connect"}
		log.Printf("connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		clients[c] = true
	}

	for {
		for client := range clients {
			_, message, err := client.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
			log.Println("debug")
		}
	}
}
