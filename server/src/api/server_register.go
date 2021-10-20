package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

// ServerRegister sshでマウント
func ServerRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "GET" {
		_, _ = fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	print(r.Header.Get("X-Forwarded-For"))

}
