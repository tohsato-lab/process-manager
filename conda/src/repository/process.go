package repository

import (
	"github.com/jmoiron/sqlx"
)

func SetProcess(db *sqlx.DB, processID string, targetFilename string, envName string) error {
	_, err := db.NamedExec(
		`INSERT INTO calc_process_table (id, target_file, env_name) VALUES (:id, :target_file, :env_name)`,
		map[string]interface{}{
			"id":          processID,
			"target_file": targetFilename,
			"env_name":    envName,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
