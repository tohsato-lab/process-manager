package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./api"
)

func main() {
	if err := os.Mkdir("dataset", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("programs", 0777); err != nil {
		fmt.Println(err)
	}
	if err := os.Mkdir("logs", 0777); err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/upload", api.UploadHander)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
