package api

import (
	"encoding/json"
	"net/http"

	"../utils"
)

// EnvInfoHandler Env一覧取得
func EnvInfoHandler(w http.ResponseWriter, r *http.Request) {

	envs := utils.GetCondaEnv()

	// return
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json, _ := json.Marshal(envs)
	w.Write(json)
}
