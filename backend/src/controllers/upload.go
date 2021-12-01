package controllers

import (
	"backend/repository"
	"backend/utils"
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
	return fileHeader.Filename, nil
}

func ServerInfo(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	calcServers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	var condaEnvs [][]byte
	for _, server := range calcServers {
		requestHTTP, err := utils.RequestHTTP(
			"GET", "http://"+server.IP+":"+server.Port+"/conda_env",
			5*time.Second,
		)
		if err != nil {
			return
		}

	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	ip := r.FormValue("ip")
	// env := r.FormValue("conda_env")
	// comment := r.FormValue("comment")
	execCount := r.FormValue("exec_count")
	// targetFile := r.FormValue("target_file")

	if num, err := strconv.Atoi(execCount); num <= 0 {
		http.Error(w, "実行回数が0以下です", http.StatusBadGateway)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	filename, err := putFile(r.FormFile("file"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	calcServers, err := repository.GetCalcServer(db, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	log.Println(calcServers)
	utils.SendFile(filename, "http://"+ip+":5984/upload")

}
