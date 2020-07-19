package api

import (
	"log"
	"net/http"

	"../utils"
	"github.com/gorilla/websocket"
)

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

// WebSocket サーバーにつなぎにいくクライアント
var clients = make(map[*websocket.Conn]bool)

// ProcessStatusHandle プロセス一覧をリアルタイムで返す
func ProcessStatusHandle(w http.ResponseWriter, r *http.Request) {
	// websocket の状態を更新

	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}

	// websocket を閉じる
	defer websocket.Close()

	clients[websocket] = true

	for {
		// メッセージ受け取り
		process := <-utils.BroadcastProcess

		// クライアントの数だけループ
		for client := range clients {
			//　書き込む
			err := client.WriteJSON(process)
			if err != nil {
				log.Printf("error occurred while writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
