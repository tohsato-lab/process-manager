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
	Pipe chan map[string]string
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

		var commands []map[string]string
		if err := json.Unmarshal(message, &commands); err != nil {
			log.Println(err)
			continue
		}
		for _, command := range commands {
			process, err := repository.GetProcess(c.DB, command["ID"])
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(process)

			switch command["status"] {
			case "running":
				go func() {
					log.Println("exec process")
					if err := repository.UpdateProcessStatus(c.DB, command["ID"], command["status"]); err != nil {
						log.Println(err)
						return
					}

					c.Pipe <- map[string]string{"ID": process.ID, "status": "running"}
					status, err := execute(c.DB, process.ID, process.TargetFile, process.EnvName)
					if err != nil {
						log.Println(err)
						return
					}
					if err := repository.UpdateProcessStatus(c.DB, command["ID"], status); err != nil {
						log.Println(err)
						return
					}
					log.Println("exec done")
					c.Pipe <- map[string]string{"ID": process.ID, "status": status}
				}()
			case "kill":
				log.Println("kill process")
				status, err := killCMD(c.DB, process.ID)
				if err != nil {
					log.Println(err)
					return
				}
				if err := repository.UpdateProcessStatus(c.DB, command["ID"], status); err != nil {
					log.Println(err)
					return
				}
				c.Pipe <- map[string]string{"ID": process.ID, "status": status}
			case "sync":
				c.Pipe <- map[string]string{"ID": process.ID, "status": process.Status}
			}
		}
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
	for {
		select {
		case result, ok := <-c.Pipe:
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println(err)
				return
			}
			contents, err := json.Marshal(result)
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
		}
	}
}
