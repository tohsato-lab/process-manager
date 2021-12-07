package modules

import (
	"github.com/jmoiron/sqlx"
	"log"

	"backend/repository"
)

func UpdateProcess(db *sqlx.DB) {
	servers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, server := range servers {
		execIDs, err := repository.CanExecProcess(db, server.IP, 1)
		if err != nil {
			log.Println(err.Error())
			return
		} else if execIDs == nil {
			log.Println("該当なし")
			return
		}
		for _, processID := range execIDs {
			ExecProcess(processID, server.IP)
		}
	}
}

func ExecProcess(processID string, serverIP string) {
	if err := connections[serverIP].WriteJSON(map[string]string{"ID": processID, "command": "running"}); err != nil {
		log.Println(err)
		return
	}
}
