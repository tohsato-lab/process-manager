package controllers

import (
	"conda/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func putFile(file multipart.File, fileHeader *multipart.FileHeader, err error) (string, error) {
	if err != nil {
		log.Println("ファイルアップロードを確認できませんでした")
		return "", err
	}
	saveFile, err := os.Create("./" + fileHeader.Filename)
	if err != nil {
		log.Println("サーバ側でファイル作成できませんでした")
		return "", err
	}
	if _, err := io.Copy(saveFile, file); err != nil {
		log.Println("アップロードしたファイルの書き込みに失敗しました")
		return "", err
	}
	if err := saveFile.Close(); err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return fileHeader.Filename, err
}

func FileUpload(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	// ip := r.FormValue("ip")
	env := r.FormValue("conda_env")
	log.Println(env)
	// comment := r.FormValue("comment")
	// execCount := r.FormValue("exec_count")
	// targetFile := r.FormValue("target_file")
	filename, err := putFile(r.FormFile("file"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	md5Data := md5.Sum([]byte(time.Now().String()))
	processID := hex.EncodeToString(md5Data[:])
	targetDIR := "/process-manager/data/" + processID + "/"
	if err := os.MkdirAll(targetDIR, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if _, err := exec.Command("sh", "-c", "unzip "+filename+" -d "+targetDIR).Output(); err != nil {
		_, _ = fmt.Fprintln(w, "ファイル解凍に失敗しました。"+err.Error())
		return
	}

	utils.RespondByte(w, http.StatusOK, []byte(`{"status":"ok"}`))
}
