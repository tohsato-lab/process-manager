package modules

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	"backend/repository"
)

const (
	writeWait = 10 * time.Second
)

type Client struct {
	DB   *sqlx.DB
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) WritePump() {
	defer func() {
		SocketCore.Unregister <- c
		if err := c.Conn.Close(); err != nil {
			return
		}
	}()
	activeProcess, err := repository.GetProcess(c.DB, false)
	if err != nil {
		log.Println(err)
		return
	}
	contents, err := json.Marshal(activeProcess)
	if err != nil {
		log.Println(err)
		return
	}
	if err := c.Conn.WriteMessage(websocket.TextMessage, contents); err != nil {
		log.Println(err)
		return
	}
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println(err)
				return
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
