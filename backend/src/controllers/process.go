package controllers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os/exec"
	"time"

	"backend/modules"
	"backend/repository"
	"backend/utils"
)

func EntryProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	log.Println("### EntryProcess")
	processName := r.FormValue("process_name")
	envName := r.FormValue("conda_env")
	serverIP := r.FormValue("server_ip")
	comment := r.FormValue("comment")
	args := r.FormValue("args")
	var processIDs []string
	if err := json.Unmarshal([]byte(r.FormValue("process_ids")), &processIDs); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	for _, processID := range processIDs {
		if err := repository.SetProcess(db, processID, processName, envName, serverIP, comment, args); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
	if err := modules.UpdateProcess(db); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func KillProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	log.Println("### KillProcess")
	processID := r.FormValue("process_id")
	serverIP, err := repository.GetProcessServerIP(db, processID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if err := modules.KillProcess(db, processID, serverIP); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func InTrashProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	processID := r.FormValue("process_id")
	if err := modules.TrashProcess(db, processID); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func InTrashAllProcess(w http.ResponseWriter, _ *http.Request, db *sqlx.DB) {
	trashProcess, err := repository.GetProcesses(db, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	contents, err := json.Marshal(trashProcess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)
}

func DeleteProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	log.Println("### DeleteProcess")
	processID := r.FormValue("process_id")
	serverIP, err := repository.GetProcessServerIP(db, processID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	server, err := repository.GetCalcServer(db, serverIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if requestHTTP, err := utils.RequestHTTP(
		"DELETE", "http://"+server.IP+":"+server.Port+"/delete?process_id="+processID, 5*time.Second,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	} else {
		var res map[string]string
		if err := json.Unmarshal(requestHTTP, &res); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if res["status"] != "ok" {
			http.Error(w, res["status"], http.StatusBadGateway)
			return
		}
	}
	if err := repository.DeleteProcess(db, processID); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	cmd := "rm -rf ../../log/" + processID + "/"
	log.Println(cmd)
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	InTrashAllProcess(w, r, db)
}
