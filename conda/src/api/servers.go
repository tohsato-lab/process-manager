package api

import (
	"conda/modules"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// Servers サーバーをDBに登録
func Servers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		if r.FormValue("mode") == "add" {
			fmt.Println("add")
			ip := r.FormValue("ip")
			port := r.FormValue("port")
			modules.RegisterServer(db, ip, port)
		} else {
			fmt.Println("delete")
			ip := r.FormValue("ip")
			modules.DeleteServer(db, ip)
		}
	default:
		fmt.Println("unknow method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	servers := modules.GetServers(db)
	fmt.Println(servers)

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	content, err := json.Marshal(servers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(content); err != nil {
		fmt.Println(err)
	}

}
