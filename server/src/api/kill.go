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
	err := db.QueryRow("SELECT pid FROM process_table WHERE id = ?", id).Scan(&pid)
	if err != nil {
		panic(err.Error())
	}

	// 親子共々kill
	exec.Command("sh", "-c", "kill `ps ho pid --ppid="+strconv.Itoa(pid)+"`").Run()
	exec.Command("kill", strconv.Itoa(pid)).Run()

	modules.ComplateProcess(db, id, "killed")

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   "killed",
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
