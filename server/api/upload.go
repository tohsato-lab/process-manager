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
	"strconv"
	"strings"
	"time"

	"../modules"
	"../utils"
)

type response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// UploadHandler ファイルアップロードハンドラー
func UploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

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

	// get use vram
	vram, e := strconv.ParseFloat(r.FormValue("vram"), 32)
	if e != nil {
		fmt.Fprintln(w, "使用VRAM容量を確認出来ませんでした。")
		return
	}

	// target filename
	md5 := md5.Sum([]byte(time.Now().String()))
	targetFileID := hex.EncodeToString(md5[:])

	// unzip
	utils.Unzip("./"+uploadedFileName, "../data/programs/"+targetFileID)
	if err := os.Rename("./"+uploadedFileName, "../data/programs/"+targetFileID+"/"+uploadedFileName); err != nil {
		fmt.Println(err)
	}

	// regist proceess
	modules.RegistProcess(db, utils.Process{
		ID:       targetFileID,
		UseVram:  float32(vram),
		Status:   "ready",
		Filename: strings.Split(uploadedFileName, ".")[0],
	})

	// update process
	modules.UpdataAllProcess(db)

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   "success",
	}
	json, _ := json.Marshal(response)
	w.Write(json)
}
