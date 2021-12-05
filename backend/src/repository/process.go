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

func CanExecProcess(db *sqlx.DB, serverIP string, limit int) ([]string, error) {
	var numRunning int
	if err := db.Get(
		&numRunning, `SELECT COUNT(*) FROM process_table WHERE status='running' AND server_ip=?`, serverIP,
	); err != nil {
		return nil, err
	}
	if numRunning >= limit {
		return nil, nil
	}

	var canExecID []string
	if err := db.Select(
		&canExecID, `SELECT id FROM process_table 
						   WHERE status='ready' AND in_trash = false AND server_ip=?
		                   ORDER BY upload_date LIMIT ?`, serverIP, limit,
	); err != nil {
		return nil, err
	}
	return canExecID, nil
}
