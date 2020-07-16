package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"../utils"
)

type response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// UploadHander ファイルアップロードハンドラー
func UploadHander(w http.ResponseWriter, r *http.Request) {
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
	_, e = io.Copy(saveFile, file)
	if e != nil {
		fmt.Println(e)
		fmt.Println("アップロードしたファイルの書き込みに失敗しました。")
		os.Exit(1)
	}

	// unzip
	utils.Unzip("./"+uploadedFileName, "../programs")

	// launch
	go Launch("../programs/test/")

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "ok",
		Data:   "test",
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
