package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"./api"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "golang:golang@/process_manager_db?parseTime=true")
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
	http.HandleFunc("/env_info", func(w http.ResponseWriter, r *http.Request) {
		api.EnvInfoHandler(w, r)
	})
	http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
		api.WebSocketHandle(w, r, db)
	})
	/*
		http.HandleFunc("/gpu_status", func(w http.ResponseWriter, r *http.Request) {
			api.GPUSstatus(w, r)
		})
	*/
	http.HandleFunc("/programs/", func(w http.ResponseWriter, r *http.Request) {
		api.Explorer(w, r, db)
	})
	go api.WebSocketKernel()

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":5983", nil))
}
