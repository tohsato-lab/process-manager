package controllers

import (
	"conda/repository"
	"conda/utils"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func putFile(file multipart.File, fileHeader *multipart.FileHeader, directoryName string) (string, error) {
	saveFile, err := os.Create(directoryName + fileHeader.Filename)
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

	envName := r.FormValue("conda_env")
	execCount, _ := strconv.Atoi(r.FormValue("exec_count"))
	targetFile := r.FormValue("target_file")
	execCount, err := strconv.Atoi(r.FormValue("exec_count"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if execCount < 1 {
		http.Error(w, "実行回数が1未満です", http.StatusBadGateway)
		return
	}

	var processIDs []string

	for i := 0; i < execCount; i++ {
		md5Data := md5.Sum([]byte(time.Now().String()))
		processID := hex.EncodeToString(md5Data[:])
		targetDIR := "/process-manager/data/" + processID + "/"
		if err := os.MkdirAll(targetDIR, 0755); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		filename, err := putFile(file, fileHeader, targetDIR)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		cmd := "unzip " + targetDIR + filename + " -d " + targetDIR
		if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if err := repository.SetProcess(db, processID, targetFile, envName); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		processIDs = append(processIDs, processID)
	}

	contents, err := json.Marshal(processIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)

}