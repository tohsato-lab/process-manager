package api

import (
	"database/sql"
	"fmt"
	"github.com/janberktold/sse"
	"net/http"
	"process-manager-server/modules"
	"process-manager-server/utils"
)

// ProcessStatus プロセスの状態を配信
func ProcessStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	conn, _ := sse.Upgrade(w, r)
	if err := conn.WriteJson(modules.GetAllProcess(db)); err != nil {
		fmt.Println(err)
		conn.Close()
	}

	for {
		// メッセージ受け取り
		process := <-utils.BroadcastProcess
		if err := conn.WriteJson(process); err != nil {
			fmt.Println(err)
			conn.Close()
		}
	}
}
