package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"./api"
	"./utils"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "golang:golang@/process_manager_db?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	utils.GetCondaEnv()

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHandler(w, r, db)
	})
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		api.KillHandler(w, r, db)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteHandler(w, r, db)
	})
	http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
		api.WebSocketHandle(w, r, db)
	})
	go api.WebSocketKernel()

	fmt.Println("server start")
	log.Fatal(http.ListenAndServe(":5983", nil))
}
