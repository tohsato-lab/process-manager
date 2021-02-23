package api

import (
	"encoding/json"
	"net/http"

	"../utils"
)

// EnvInfoHandler Env一覧取得
func EnvInfoHandler(w http.ResponseWriter, _ *http.Request) {

	envs := utils.GetCondaEnv()

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	jsonData, _ := json.Marshal(envs)
	if _, err := w.Write(jsonData); err != nil {
		return
	}
}
