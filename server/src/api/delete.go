package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// DeleteHandler process命令実行
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	id := r.FormValue("id")

	// pid取得
	dbDelete, err := db.Prepare("DELETE FROM process_table WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer dbDelete.Close()

	result, err := dbDelete.Exec(id)
	if err != nil {
		panic(err.Error())
	}

	rowsAffect, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect)

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   "deleted",
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
