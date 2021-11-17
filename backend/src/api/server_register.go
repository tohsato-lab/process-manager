package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"process-manager-server/modules"
)

// ServerRegister サーバーをDBに登録
func ServerRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "POST" {
		_, _ = fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}

	// get ip
	ip := r.FormValue("ip")

	// get port
	port := r.FormValue("port")

	modules.RegisterServer(db, ip, port)

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
	}
	jsonData, _ := json.Marshal(response)
	if _, err := w.Write(jsonData); err != nil {
		fmt.Println(err)
	}
}
