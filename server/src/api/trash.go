package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"process-manager-server/modules"
)

// TrashHandler ゴミ箱に移動
func TrashHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.FormValue("id")

	dbStatus := ""
	var inTrash bool
	err := db.QueryRow(
		"SELECT status, in_trash FROM main_processes WHERE id = ?", id,
	).Scan(&dbStatus, &inTrash)
	if err != nil {
		fmt.Println(err)
	}

	status := "in_trash"
	if dbStatus != "running" {
		modules.TrashProcess(db, id)
	} else {
		status = "not_in_trash"
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
