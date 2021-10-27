package api

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"log"

	"net/http"
	"process-manager-server/modules"
	"process-manager-server/utils"
)

// WebSocket サーバーにつなぎにいくクライアント
var clients = make(map[*websocket.Conn]bool)

// WebSocket 更新用
// request origin not allowed by Upgrader.CheckOrigin エラーを突破
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Add this lines
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ProcessStatus ソケット接続
func ProcessStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// websocket の状態を更新
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	clients[socket] = true
	utils.BroadcastProcesses <- modules.GetProcesses(db)
}

// ProcessStatusKernel WebSocketで情報を投げる
func ProcessStatusKernel() {
	for {
		// メッセージ受け取り
		process := <-utils.BroadcastProcesses
		// クライアントの数だけループ
		for client := range clients {
			//　書き込む
			err := client.WriteJSON(process)
			if err != nil {
				log.Printf("error occurred while writing message to client: %v", err)
				if err := client.Close(); err != nil {
					return
				}
				delete(clients, client)
			}
		}
	}
}

