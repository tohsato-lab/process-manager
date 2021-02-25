package api

import (
	"database/sql"
	"fmt"
	"github.com/janberktold/sse"
	"net/http"
	"process-manager-server/modules"
	"process-manager-server/utils"
	"time"
)

// ProcessStatus プロセスの状態を配信
func ProcessStatus(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	upGrader := sse.Upgrader{
		RetryTime: 1 * time.Second,
	}

	conn, _ := upGrader.Upgrade(w, r)
	if err := conn.WriteJson(modules.GetAllProcess(db)); err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	for {
		// メッセージ受け取り
		process := <-utils.BroadcastProcess
		fmt.Println("## send process status")
		fmt.Println(process[0])
		if err := conn.WriteJson(process); err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
	}
}
