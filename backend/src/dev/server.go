package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{} // use default options

func connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
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
	// http.HandleFunc("/connect", connect)
	// log.Fatal(http.ListenAndServe(":8100", nil))

	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	abc := make(map[chan int]bool)
	abc[a] = true
	abc[b] = true
	abc[c] = true

	go func() { for { a <- 0 } }()
	go func() { for { <-c } }()
	go func() { for { b <- 0 } }()

	for i := 0; i < 10; i++ {
		for ints := range abc {
			select {
			case <-ints:
				println("debug")
			}
		}
	}

}
