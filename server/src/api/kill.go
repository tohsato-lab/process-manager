package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"

	"../modules"
)

// KillHander kill命令実行
func KillHander(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := r.FormValue("id")

	// pid取得
	var pid = 0
	err := db.QueryRow("SELECT IFNULL(pid, 0) FROM process_table WHERE id = ?", id).Scan(&pid)
	if err != nil {
		panic(err.Error())
	}

	status := "killed"
	if pid != 0 {
		// 親子共々kill
		if err := exec.Command("sh", "-c", "kill `ps ho pid --ppid="+strconv.Itoa(pid)+"`").Run(); err != nil {
			panic(err.Error())
		}
		if err := exec.Command("kill", strconv.Itoa(pid)).Run(); err != nil {
			panic(err.Error())
		}
		// DBアップデート
		modules.ComplateProcess(db, id, status)
	} else {
		status = "not kill"
	}

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   status,
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
