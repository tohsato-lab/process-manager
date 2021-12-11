package controllers

import (
	"backend/repository"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"path/filepath"
)

type SpaHandler struct {
	StaticPath string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(h.StaticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}

func ProcessLog(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	processID := r.FormValue("process_id")
	process, err := repository.GetProcess(db, processID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	var hostName string
	if process.Status == "ready" || process.Status == "running" {
		server, err := repository.GetCalcServer(db, process.ServerIP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		hostName = server.IP + ":" + server.Port
	} else {
		hostName = r.Host
	}
	http.Redirect(w, r, "http://"+hostName+"/data/"+processID, http.StatusMovedPermanently)
}
