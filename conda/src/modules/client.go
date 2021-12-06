package modules

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		log.Println("disconnect")
		c.Hub.unregister <- c
		if err := c.Conn.Close(); err != nil {
			return
		}
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println(string(message))
		c.Send <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("destroy Write Pump")
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			return
		}
	}()
	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("hi.")); err != nil {
		log.Println(err)
		return
	}
	for {
		log.Println("wait")
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				return
			}
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Println(err)
					return
				}
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Println(err)
				return
			}
			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			log.Println("ping")
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-time.After(2 * time.Second):
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Println(err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("hi?")); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
