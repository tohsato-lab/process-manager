package controllers

import (
	"github.com/jmoiron/sqlx"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	log.Println("done.")
	return fileHeader.Filename, nil
}

func FileUpload(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

	log.Println("Upload")
	// ip := r.FormValue("ip")
	// env := r.FormValue("conda_env")
	// comment := r.FormValue("comment")
	// execCount := r.FormValue("exec_count")
	// targetFile := r.FormValue("target_file")

	log.Println(r.FormFile("file"))
	log.Println("debug")
	_, err := putFile(r.FormFile("file"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

}
