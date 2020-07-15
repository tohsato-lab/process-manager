package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/upload", uploadHander)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GET hello!\n")
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "POST hello!\n")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Method not allowed.\n")
	}
}

func uploadHander(w http.ResponseWriter, r *http.Request) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "POST" {
		fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var e error
	var uploadedFileName string
	// POSTされたファイルデータを取得する
	file, fileHeader, e = r.FormFile("file")
	fmt.Printf("%T\n", file)
	if e != nil {
		fmt.Fprintln(w, "ファイルアップロードを確認できませんでした。")
		return
	}
	uploadedFileName = fileHeader.Filename
	var saveFile *os.File
	saveFile, e = os.Create("./" + uploadedFileName)
	if e != nil {
		fmt.Fprintln(w, "サーバ側でファイル確保できませんでした。")
		return
	}
	defer saveFile.Close()
	defer file.Close()
	size, e := io.Copy(saveFile, file)
	if e != nil {
		fmt.Println(e)
		fmt.Println("アップロードしたファイルの書き込みに失敗しました。")
		os.Exit(1)
	}
	fmt.Println("書き込んだByte数=>")
	fmt.Println(size)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := Response{
		Status: "ok",
		Data:   "test",
	}
	json, _ := json.Marshal(response)

	w.Write(json)
}
