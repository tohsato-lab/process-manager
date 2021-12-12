package controllers

import (
	"conda/utils"
	"encoding/json"
	"net/http"
)

func EnvInfo(w http.ResponseWriter, _ *http.Request) {
	condaEnv, err := utils.GetCondaEnv()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	contents, err := json.Marshal(condaEnv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)
}