package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"process-manager-server/modules"
	"strconv"
)

// KillHandler kill命令実行
func KillHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := r.FormValue("id")

	// pid取得
	pid := 0
	dbStatus := ""
	err := db.QueryRow("SELECT IFNULL(pid, 0), status FROM main_processes WHERE id = ?", id).Scan(&pid, &dbStatus)
	if err != nil {
		fmt.Println(err)
	}

	status := "killed"
	if pid != 0 {
		// 親子共々 kill
		if err := exec.Command("sh", "-c", "kill `ps ho pid --ppid="+strconv.Itoa(pid)+"`").Run(); err != nil {
			fmt.Println(err)
			status = "not kill"
		}
	} else {
		status = "not kill"
	}

	//何らかの理由でステータスが更新されていなかった場合はkilledに更新
	if dbStatus == "running" && status == "not kill" {
		// DBアップデート
		modules.CompleteProcess(db, id, "killed")
	}

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   status,
	}
	jsonData, _ := json.Marshal(response)
	if _, err := w.Write(jsonData); err != nil {
		fmt.Println(err)
	}
}
