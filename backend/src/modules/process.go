package modules

import (
	"backend/repository"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
)

func syncProcess(db *sqlx.DB) error {
	activeProcess, err := repository.GetProcesses(db, false)
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
	commands := []map[string]string{{"ID": processID, "status": "running"}}
	if err := connections[serverIP].WriteJSON(commands); err != nil {
		return err
	}
	return nil
}

func KillProcess(db *sqlx.DB, processID string, serverIP string) error {
	if err := repository.UpdateProcessStatus(db, processID, "killed"); err != nil {
		return err
	}
	commands := []map[string]string{{"ID": processID, "status": "kill"}}
	if err := connections[serverIP].WriteJSON(commands); err != nil {
		return err
	}
	return nil
}

func TrashProcess(db *sqlx.DB, processID string) error {
	log.Println("trash")
	if err := repository.UpdateProcessTrash(db, processID); err != nil {
		return err
	}
	if err := syncProcess(db); err != nil {
		return err
	}
	return nil
}
