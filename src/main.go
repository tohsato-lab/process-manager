package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"./api"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := os.Mkdir("../dataset", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("../programs", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("../histories", 0777); err != nil {
		fmt.Println(err)
	}
	db, err := sql.Open("mysql", "golang:golang@/process_manager_db?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHander(w, r, db)
	})
	http.Handle("/histories/", http.StripPrefix("/histories/", http.FileServer(http.Dir("../histories"))))
	http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
		api.WebSocketHandle(w, r, db)
	})
	go api.WebSocketKernel()

	log.Fatal(http.ListenAndServeTLS(":5983", "../ssl/cert.pem", "../ssl/key.pem", nil))
}
