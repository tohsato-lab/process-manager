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
			log.Println(processID)
			ExecProcess(db, processID)
		}
	}

}

func ExecProcess(db *sqlx.DB, processID string) {
	log.Println("ExecProcess")
}
