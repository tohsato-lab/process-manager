package modules

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"

	"backend/repository"
)

func syncProcess(db *sqlx.DB) error {
	activeProcess, err := repository.GetActiveProcess(db)
	if err != nil {
		return err
	}
	processes, err := json.Marshal(activeProcess)
	if err != nil {
		return err
	}
	SocketCore.Broadcast <- processes
	return nil
}

func UpdateProcess(db *sqlx.DB) error {
	servers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		return err
	}
	for _, server := range servers {
		execIDs, err := repository.CanExecProcess(db, server.IP, 1)
		if err != nil {
			return err
		} else if execIDs == nil {
			log.Println("該当なし")
			if err := syncProcess(db); err != nil {
				return err
			}
			continue
		}
		for _, processID := range execIDs {
			if err := ExecProcess(db, processID, server.IP); err != nil {
				return err
			}
		}
	}
	return nil
}

func ExecProcess(db *sqlx.DB, processID string, serverIP string) error {
	if err := repository.UpdateProcessStatus(db, processID, "running"); err != nil {
		return err
	}
	if err := connections[serverIP].WriteJSON(map[string]string{"ID": processID, "status": "running"}); err != nil {
		return err
	}
	return nil
}

func KillProcess(db *sqlx.DB, processID string, serverIP string) error {
	if err := repository.UpdateProcessStatus(db, processID, "killed"); err != nil {
		return err
	}
	if err := connections[serverIP].WriteJSON(map[string]string{"ID": processID, "status": "kill"}); err != nil {
		return err
	}
	return nil
}
