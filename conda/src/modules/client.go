package modules

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	"conda/repository"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	DB   *sqlx.DB
	Conn *websocket.Conn
	Pipe chan string
}

var Clients = make(map[*Client]bool)

func (c *Client) ReadPump() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println(string(message))

		command := map[string]string{}
		if err := json.Unmarshal(message, &command); err != nil {
			log.Println(err)
			continue
		}
		process, err := repository.GetProcess(c.DB, command["ID"])
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(process)

		switch command["command"] {
		case "running":
			go func() {
				status := Execute(c.DB, process.ID, process.TargetFile, process.EnvName)
				if err := repository.UpdateProcessStatus(c.DB, process.ID, status); err != nil {
					log.Println(err)
					return
				}
				c.Pipe <- process.ID
			}()
		case "kill":
		case "delete":

		}
		c.Pipe <- process.ID
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.Conn.Close(); err != nil {
			log.Println(err)
			return
		}
		delete(Clients, c)
		close(c.Pipe)
		log.Println("destroyed socket")
	}()
	if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("hi.")); err != nil {
		log.Println(err)
		return
	}
	for {
		select {
		case processID, ok := <-c.Pipe:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println(err)
				return
			}
			process, err := repository.GetProcess(c.DB, processID)
			if err != nil {
				log.Println(err)
				return
			}
			contents, err := json.Marshal(map[string]string{"ID": processID, "status": process.Status})
			if err != nil {
				log.Println(err)
				return
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.TextMessage, contents); err != nil {
				log.Println(err)
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-time.After(2 * time.Second):
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("hi?")); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
