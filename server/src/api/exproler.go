package api

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Explorer ディレクトリの閲覧
func Explorer(w http.ResponseWriter, r *http.Request, _ *sql.DB) {
	println(r.URL.Path[1:])
	files, err := ioutil.ReadDir("../../data/" + r.URL.Path[1:])
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	http.ServeFile(w, r, "../../data/"+r.URL.Path[1:])
}
