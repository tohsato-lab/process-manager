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
	"os/exec"
	"strconv"
	"strings"
	"time"

	"process-manager-server/modules"
	"process-manager-server/utils"
)

type response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// UploadHandler ファイルアップロードハンドラー
func UploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	println("upload handler")

	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "POST" {
		if _, err := fmt.Fprintln(w, "許可したメソッドとはことなります。"); err != nil {
			return
		}
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
		if _, err := fmt.Fprintln(w, "ファイルアップロードを確認できませんでした。"); err != nil {
			return
		}
		return
	}
	uploadedFileName = fileHeader.Filename
	saveFile, e = os.Create("./" + uploadedFileName)
	if e != nil {
		if _, err := fmt.Fprintln(w, "サーバ側でファイル確保できませんでした。"); err != nil {
			return
		}
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
		if _, err := fmt.Fprintln(w, "使用VRAM容量を確認出来ませんでした。"); err != nil {
			return
		}
		return
	}

	// get use env
	env := r.FormValue("env")

	// get use target
	target := r.FormValue("target")

	// target filename
	md5Data := md5.Sum([]byte(time.Now().String()))
	targetFileID := hex.EncodeToString(md5Data[:])

	// unzip
	targetDIR := "../../data/programs/" + targetFileID + "/"
	if err := os.Mkdir(targetDIR, 0777); err != nil {
		panic(err)
	}
	if err := os.Rename("./"+uploadedFileName, targetDIR+uploadedFileName); err != nil {
		panic(err)
	}
	if _, err := exec.Command("sh", "-c", "unzip "+targetDIR+uploadedFileName+" -d "+targetDIR).Output(); err != nil {
		panic(err)
	}

	// regist proceess
	modules.RegisterProcess(db, utils.Process{
		ID:         targetFileID,
		UseVram:    float32(vram),
		Status:     "ready",
		Filename:   strings.Split(uploadedFileName, ".")[0],
		TargetFile: target,
		EnvName:    env,
	})

	// update process
	modules.UpdateAllProcess(db)
	println("アップロード完了")

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response := response{
		Status: "200",
		Data:   "success",
	}
	jsonData, _ := json.Marshal(response)
	w.Write(jsonData);
}
