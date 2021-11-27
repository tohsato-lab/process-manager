package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"backend/modules"
)

func JoinServer(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case http.MethodPost:
		ip := r.FormValue("ip")
		port := r.FormValue("port")
		errConn := modules.Connection(ip, port)
		if errConn != nil {
			http.Error(w, errConn.Error(), http.StatusBadGateway)
			return
		}
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodGet:
		servers := modules.GetServers(db)
		content, _ := json.Marshal(servers)
		_, err := w.Write(content)
		if err != nil {
			log.Println(err)
			return
		}
	case http.MethodOptions:
		return
	default:
		http.Error(w, "unKnow Method", http.StatusMethodNotAllowed)
		return
	}
}
