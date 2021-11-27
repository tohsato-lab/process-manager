package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"conda/utils"
)

// Explorer ディレクトリの閲覧
func Explorer(w http.ResponseWriter, r *http.Request, _ *sql.DB) {
	println(r.URL.Path[1:])
	var info []utils.DirectoryInfo
	files, err := ioutil.ReadDir("../../data/" + r.URL.Path[1:])
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
		info = append(info, utils.DirectoryInfo{Name: file.Name(), IsDir: file.IsDir()})
	}
	fmt.Println(info)
	http.ServeFile(w, r, "../../data/"+r.URL.Path[1:])
}
