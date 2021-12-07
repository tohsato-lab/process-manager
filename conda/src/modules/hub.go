package modules

import "log"

type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			log.Println("Register")
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("disconnected")
			if _, ok := h.clients[client]; ok {
				log.Println("disconnected")
				delete(h.clients, client)
				close(client.Pipe)
			}
		}
	}
}
