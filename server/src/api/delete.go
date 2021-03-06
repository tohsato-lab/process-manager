package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"process-manager-server/modules"
)

// DeleteHandler 完全削除
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := r.FormValue("id")

	dbStatus := ""
	var inTrash bool
	err := db.QueryRow("SELECT status, in_trash FROM main_processes WHERE id = ?", id).Scan(&dbStatus, &inTrash)
	if err != nil {
		fmt.Println(err)
	}

	status := ""
	if dbStatus != "running" && inTrash {
		targetDIR := "../../data/programs/" + id + "/"
		if _, err := exec.Command("sh", "-c", "rm -rf "+targetDIR).Output(); err != nil {
			_, _ = fmt.Fprintln(w, "ファイル削除に失敗しました。"+err.Error())
			return
		}
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
