package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"process-manager-server/modules"
)

// DeleteHandler process命令実行
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := r.FormValue("id")

	dbStatus := ""
	err := db.QueryRow("SELECT status FROM process_table WHERE id = ?", id).Scan(&dbStatus)
	if err != nil {
		fmt.Println(err)
	}

	status := ""
	if dbStatus != "running" {
		modules.DeleteProcess(db, id)
		status = "deleted"
	} else {
		status = "not delete"
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
