package api

import (
	"database/sql"
	"net/http"

	"backend/modules"
)

func JoinServer(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodPost:
		ip := r.FormValue("ip")
		port := r.FormValue("port")
		err := modules.Connection(ip, port)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

	case http.MethodGet:
	default:
		http.Error(w, "unKnow Method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte("ok"))
}
