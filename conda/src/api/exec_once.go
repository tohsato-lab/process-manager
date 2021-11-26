package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"process-manager-server/modules"
	"process-manager-server/utils"
)

// ExecOnce プロセスを登録し実行
func ExecOnce(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != http.MethodPost {
		fmt.Println("unknow method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get use env
	targetFileID := r.FormValue("targetFileID")

	// get use env
	env := r.FormValue("env")

	// get use target
	target := r.FormValue("target")

	// get use target
	comment := r.FormValue("comment")

	// register process
	modules.RegisterProcess(db, utils.Process{
		ID:         targetFileID,
		Status:     "ready",
		TargetFile: target,
		EnvName:    env,
		Comment:    comment,
		// HomeIP:     "localhost",
		InTrash:    false,
	})

}
