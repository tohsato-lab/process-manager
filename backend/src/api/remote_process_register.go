package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"process-manager-server/modules"
	"process-manager-server/utils"
)

// RemoteProcessRegister プロセスをリモートから登録
func RemoteProcessRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "GET" {
		_, _ = fmt.Fprintln(w, "許可したメソッドとはことなります。")
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
