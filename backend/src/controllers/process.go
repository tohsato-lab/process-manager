package controllers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	"backend/modules"
	"backend/repository"
	"backend/utils"
)

func EntryProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	processName := r.FormValue("process_name")
	envName := r.FormValue("conda_env")
	serverIP := r.FormValue("server_ip")
	comment := r.FormValue("comment")

	var processIDs []string
	if err := json.Unmarshal([]byte(r.FormValue("process_ids")), &processIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	for _, processID := range processIDs {
		if err := repository.SetProcess(db, processID, processName, envName, serverIP, comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}

	if err := modules.UpdateProcess(db); err != nil {
		log.Println(err)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))

}

func KillProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	processID := r.FormValue("process_id")
	serverIP, err := repository.GetProcessServerIP(db, processID)
	if err != nil {
		log.Println(err)
		return
	}
	if err := modules.KillProcess(db, processID, serverIP); err != nil {
		log.Println(err)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func DeleteProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

}
