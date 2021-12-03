package controllers

import (
	"backend/repository"
	"backend/utils"
	"encoding/json"
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

func ServerInfo(w http.ResponseWriter, _ *http.Request, db *sqlx.DB) {
	calcServers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	execInfos := map[string][]string{}
	for _, server := range calcServers {
		requestHTTP, err := utils.RequestHTTP(
			"GET", "http://"+server.IP+":"+server.Port+"/conda", 5*time.Second,
		)
		if err != nil {
			return
		}
		var envs []string
		if err := json.Unmarshal(requestHTTP, &envs); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		execInfos[server.IP] = envs
	}
	contents, err := json.Marshal(execInfos)
	log.Println(string(contents))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)
}

func Upload(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	ip := r.FormValue("ip")
	env := r.FormValue("conda_env")
	comment := r.FormValue("comment")
	execCount := r.FormValue("exec_count")
	targetFile := r.FormValue("target_file")

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
	form := map[string]string{
		"ip":          ip,
		"conda_env":   env,
		"comment":     comment,
		"exec_count":  execCount,
		"target_file": targetFile,
	}
	body, err := utils.SendFile(filename, "http://"+ip+":"+calcServers.Port+"/upload", form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	log.Println(string(body))

}
