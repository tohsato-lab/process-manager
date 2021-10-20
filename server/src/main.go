package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"process-manager-server/api"
)

func main() {
	db, err := sql.Open("mysql", "golang:process_manager@/process_manager_db?parseTime=true")
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

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":5983", nil))
}
