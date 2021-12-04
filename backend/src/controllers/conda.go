package controllers

import (
	"backend/repository"
	"backend/utils"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type Property struct {
	IP   string
	Port string
	Envs []string
}

func Envs(w http.ResponseWriter, _ *http.Request, db *sqlx.DB) {
	calcServers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if len(calcServers) == 0 {
		http.Error(w, "利用できるサーバーが見つかりません", http.StatusBadGateway)
		return
	}
	execInfos := map[string]Property{}
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
		execInfos[server.IP] = Property{
			IP:   server.IP,
			Port: server.Port,
			Envs: envs,
		}
	}
	contents, err := json.Marshal(execInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	utils.RespondByte(w, http.StatusOK, contents)
}
