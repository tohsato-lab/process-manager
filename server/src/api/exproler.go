package api

import (
	"database/sql"
	"net/http"
)

// Explorer ディレクトリの閲覧
func Explorer(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	println(r.URL.Path[1:])
	http.ServeFile(w, r, "../../data/"+r.URL.Path[1:])
}
