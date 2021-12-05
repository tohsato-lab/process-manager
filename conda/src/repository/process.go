package repository

import (
	"github.com/jmoiron/sqlx"
	"time"
)

func SetProcess(db *sqlx.DB, processID string, targetFilename string, envName string) error {
	_, err := db.NamedExec(
		`INSERT INTO calc_process_table (id, target_file, env_name, upload_date) 
			   VALUES (:id,:target_file,:env_name, :upload_date)`, map[string]interface{}{
			"id":          processID,
			"target_file": targetFilename,
			"env_name":    envName,
			"upload_date": time.Now(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
