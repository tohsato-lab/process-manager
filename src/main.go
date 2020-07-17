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

// Database データベース操作用
type Database struct {
	Val string
}

func main() {
	if err := os.Mkdir("../dataset", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("../programs", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("../logs", 0777); err != nil {
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

	log.Fatal(http.ListenAndServe(":8081", nil))
}
