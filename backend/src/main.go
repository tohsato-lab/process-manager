package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"

	"backend/api"
)

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "docker:docker@tcp(mysql_host:3306)/process_manager_db?parseTime=true"
)

func main() {
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if db.Close() != nil {
			fmt.Println(err)
		}
	}(db)

	http.HandleFunc("/join_server", func(w http.ResponseWriter, r *http.Request) {
		api.JoinServer(w, r, db)
	})

	/*
		http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
				api.UploadHandler(w, r, db)
		})
		http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
			api.ProcessStatus(w, r, db)
		})
		http.HandleFunc("/programs/", func(w http.ResponseWriter, r *http.Request) {
			api.Explorer(w, r, db)
		})
	*/

	log.Println("backend start")
	log.Fatal(http.ListenAndServe(":5983", nil))
}
