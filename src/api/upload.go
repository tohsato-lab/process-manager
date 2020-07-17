package api

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"../modules"
	"../utils"
)

type response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// UploadHander ファイルアップロードハンドラー
func UploadHander(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "POST" {
		fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	var file multipart.File
	var saveFile *os.File
	var fileHeader *multipart.FileHeader
	var e error
	var uploadedFileName string

	// get file
	file, fileHeader, e = r.FormFile("file")
	if e != nil {
		fmt.Fprintln(w, "ファイルアップロードを確認できませんでした。")
		return
	}
	uploadedFileName = fileHeader.Filename
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

	// target filename
	md5 := md5.Sum([]byte(time.Now().String()))
	targetFileID := hex.EncodeToString(md5[:])

	// unzip
	utils.Unzip("./"+uploadedFileName, "../programs/"+targetFileID)

	// regist proceess
	modules.RegistProcess(db, &modules.Process{
		ID:       targetFileID,
		UseVram:  0.0,
		Status:   "ready",
		Filename: strings.Split(uploadedFileName, ".")[0],
	})

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "ok",
		Data:   "success",
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
