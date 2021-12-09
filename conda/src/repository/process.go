package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Process struct {
	ID         string        `db:"id"`
	TargetFile string        `db:"target_file"`
	EnvName    string        `db:"env_name"`
	Status     string        `db:"status"`
	PID        sql.NullInt32 `db:"pid"`
}

func GetProcess(db *sqlx.DB, processID string) (Process, error) {
	var process Process
	if err := db.Get(&process, `SELECT * FROM calc_process_table WHERE id=?`, processID); err != nil {
		return process, err
	}
	return process, nil
}

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

func SetPID(db *sqlx.DB, processID string, pid int) error {
	_, err := db.NamedExec(`UPDATE calc_process_table SET pid=:pid WHERE id=:id`,
		map[string]interface{}{"id": processID, "pid": pid},
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProcessStatus(db *sqlx.DB, processID string, status string) error {
	_, err := db.NamedExec(`UPDATE calc_process_table SET status=:status WHERE id=:id`,
		map[string]interface{}{"id": processID, "status": status},
	)
	if err != nil {
		return err
	}
	return nil
}
