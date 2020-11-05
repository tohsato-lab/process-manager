package api

import (
	"database/sql"
	"net/http"
)

// Exproler kill命令実行
func Exproler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.FormValue("id")
	println(id)
	println(r.URL.Path[1:])
	http.ServeFile(w, r, "../../data/"+r.URL.Path[1:])
}
