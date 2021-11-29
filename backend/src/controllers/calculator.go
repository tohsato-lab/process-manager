package controllers

import (
	"backend/repository"
	"backend/utils"
	"github.com/jmoiron/sqlx"
	"net/http"

	"backend/modules"
)

func JoinServer(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
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
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}

func ServerStatus(w http.ResponseWriter, _ *http.Request, db *sqlx.DB) {
	content, err := repository.GetCalcServers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusVariantAlsoNegotiates)
		return
	}
	utils.RespondByte(w, http.StatusOK, content)
}

func DeleteServer(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	ip := r.FormValue("ip")
	if err := repository.DeleteCalcServer(db, ip); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}
