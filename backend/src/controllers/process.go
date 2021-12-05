package controllers

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

func EntryProcess(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	// env := r.FormValue("conda_env")
	// comment := r.FormValue("comment")
	// execCount := r.FormValue("exec_count")

}
