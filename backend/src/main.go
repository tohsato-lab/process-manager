package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"process-manager-server/api"
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
	defer db.Close()

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHandler(w, r, db)
	})
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		api.KillHandler(w, r, db)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteHandler(w, r, db)
	})
	http.HandleFunc("/trash", func(w http.ResponseWriter, r *http.Request) {
		api.TrashHandler(w, r, db)
	})
	http.HandleFunc("/env_info", func(w http.ResponseWriter, r *http.Request) {
		api.EnvInfoHandler(w, r)
	})
	http.HandleFunc("/host_status", func(w http.ResponseWriter, r *http.Request) {
		api.HostStatus(w, r)
	})
	http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
		api.ProcessStatus(w, r, db)
	})
	http.HandleFunc("/programs/", func(w http.ResponseWriter, r *http.Request) {
		api.Explorer(w, r, db)
	})
	go api.ProcessStatusKernel()

	// サーバー
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		api.ServerRegister(w, r, db)
	})

	fmt.Println("backend start")
	log.Fatal(http.ListenAndServe(":5983", nil))
}
