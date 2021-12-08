package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Process struct {
	ID           string `db:"id"`
	ProcessName  string `db:"process_name"`
	EnvName      string `db:"env_name"`
	ServerIP     string `db:"server_ip"`
	Comment      string `db:"comment"`
	Status       string `db:"status"`
	StartDate    string `db:"start_date"`
	CompleteDate string `db:"complete_date"`
}

func GetActiveProcess(db *sqlx.DB) ([]Process, error) {
	var activeProcess []Process
	rows, err := db.Query(
		`SELECT id, process_name, env_name, server_ip, comment, status, start_date, complete_date 
			   FROM process_table WHERE !in_trash ORDER BY upload_date DESC`,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var process Process
		var startDate sql.NullTime
		var completeDate sql.NullTime
		err := rows.Scan(
			&process.ID,
			&process.ProcessName,
			&process.EnvName,
			&process.ServerIP,
			&process.Comment,
			&process.Status,
			&startDate,
			&completeDate,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if startDate.Valid {
			process.StartDate = startDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if completeDate.Valid {
			process.CompleteDate = completeDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		activeProcess = append(activeProcess, process)
	}
	return activeProcess, nil
}

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

func UpdateProcessStatus(db *sqlx.DB, processID string, status string) error {
	_, err := db.NamedExec(`UPDATE process_table SET status=:status WHERE id=:id`,
		map[string]interface{}{"id": processID, "status": status},
	)
	if err != nil {
		return err
	}
	return nil
}

func SetStartDate(db *sqlx.DB, processID string) error {
	_, err := db.NamedExec(
		`UPDATE process_table SET start_date=:start_date WHERE id=:id AND start_date IS NULL`,
		map[string]interface{}{"id": processID, "start_date": time.Now()},
	)
	if err != nil {
		return err
	}
	return nil
}

func SetCompleteDate(db *sqlx.DB, processID string) error {
	_, err := db.NamedExec(
		`UPDATE process_table SET complete_date=:complete_date WHERE id=:id AND complete_date IS NULL`,
		map[string]interface{}{"id": processID, "complete_date": time.Now()},
	)
	if err != nil {
		return err
	}
	return nil
}
