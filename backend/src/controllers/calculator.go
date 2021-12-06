package controllers

import (
	"backend/modules"
	"backend/repository"
	"backend/utils"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func JoinServer(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.FormValue("mode") == "join" {
		ip := r.FormValue("ip")
		port := r.FormValue("port")
		if err := modules.Connection(ip, port); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if err := repository.SetCalcServer(db, ip, port); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	} else {
		ip := r.FormValue("ip")
		status := r.FormValue("mode")
		port := r.FormValue("port")
		if status == "active" {
			if err := modules.Connection(ip, port); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		} else {
			log.Println("try websocket disconnect")
			if err := modules.Disconnection(ip); err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		}
		if err := repository.UpdateCalcServerStatus(db, ip, status); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func ServerStatus(w http.ResponseWriter, _ *http.Request, db *sqlx.DB) {
	calcServers, err := repository.GetAllCalcServers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	contents, err := json.Marshal(calcServers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)
}

func DeleteServer(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	ip := r.FormValue("ip")
	if err := repository.DeleteCalcServer(db, ip); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}
