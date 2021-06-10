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

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "POST" {
		_, _ = fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	var file multipart.File
	var saveFile *os.File
	var fileHeader *multipart.FileHeader
	var uploadedFileName string

	// get use env
	env := r.FormValue("env")

	// get use target
	target := r.FormValue("target")

	// get use target
	comment := r.FormValue("comment")

	// get use target
	execCount, err := strconv.Atoi(r.FormValue("exec_count"))
	if err != nil {
		fmt.Println("Convert Error.")
	}
	if execCount <= 0 {
		_, _ = fmt.Fprintln(w, "実行回数が0です。")
		return
	}

	// get file
	file, fileHeader, err = r.FormFile("file")
	if err != nil {
		_, _ = fmt.Fprintln(w, "ファイルアップロードを確認できませんでした。")
		return
	}
	uploadedFileName = fileHeader.Filename
	saveFile, err = os.Create("./" + uploadedFileName)
	if err != nil {
		_, _ = fmt.Fprintln(w, "サーバ側でファイル作成できませんでした。")
		return
	}

	if _, err := io.Copy(saveFile, file); err != nil {
		_, _ = fmt.Fprintln(w, "アップロードしたファイルの書き込みに失敗しました。")
		return
	}
	if err := saveFile.Close(); err != nil {
		fmt.Println(err)
	}
	if err := file.Close(); err != nil {
		fmt.Println(err)
	}

	for i := 0; i < execCount; i++ {
		// target filename
		md5Data := md5.Sum([]byte(time.Now().String()))
		targetFileID := hex.EncodeToString(md5Data[:])

		// unzip
		targetDIR := "../../data/programs/" + targetFileID + "/"
		if err := os.MkdirAll(targetDIR, 0777); err != nil {
			_, _ = fmt.Fprintln(w, "ディレクトリ生成に失敗しました。"+err.Error())
			return
		}
		if _, err := exec.Command("sh", "-c", "cp "+uploadedFileName+" "+targetDIR+uploadedFileName).Output(); err != nil {
			_, _ = fmt.Fprintln(w, "ファイルコピーに失敗しました。"+err.Error())
			return
		}
		if _, err := exec.Command("sh", "-c", "unzip "+targetDIR+uploadedFileName+" -d "+targetDIR).Output(); err != nil {
			_, _ = fmt.Fprintln(w, "ファイル解凍に失敗しました。"+err.Error())
			return
		}

		// register process
		modules.RegisterProcess(db, utils.Process{
			ID:         targetFileID,
			Status:     "ready",
			Filename:   strings.Split(uploadedFileName, ".")[0],
			TargetFile: target,
			EnvName:    env,
			Comment:    comment,
			InTrash:    false,
		})
	}

	if err := os.Remove(uploadedFileName); err != nil {
		fmt.Println(err)
	}

	println("アップロード完了")

	// return
	response := response{
		Status: "200",
		Data:   "success",
	}
	jsonData, _ := json.Marshal(response)
	if _, err := w.Write(jsonData); err != nil {
		fmt.Println(err)
	}
}
