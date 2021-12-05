package repository

import (
	"github.com/jmoiron/sqlx"
	"time"
)

func SetProcess(db *sqlx.DB, processID string, processName string, envName string, IP string, comment string) error {
	_, err := db.NamedExec(
		`INSERT INTO process_table (id, process_name, env_name, server_ip, comment, upload_date) 
			   VALUES (:id, :process_name, :env_name, :server_ip, :comment, :upload_date)`, map[string]interface{}{
			"id":           processID,
			"process_name": processName,
			"env_name":     envName,
			"server_ip":    IP,
			"comment":      comment,
			"upload_date":  time.Now(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
