package repository

import (
	"github.com/jmoiron/sqlx"
)

type CalcServers struct {
	IP     string `db:"ip"`
	Port   string `db:"port"`
	Status string `db:"status"`
}

func GetAllCalcServers(db *sqlx.DB) ([]CalcServers, error) {
	var calcServers []CalcServers
	if err := db.Select(&calcServers, `SELECT * FROM calc_server_table`); err != nil {
		return nil, err
	}
	return calcServers, nil
}

func GetActiveCalcServers(db *sqlx.DB) ([]CalcServers, error) {
	var calcServers []CalcServers
	if err := db.Select(&calcServers, `SELECT * FROM calc_server_table WHERE status='active'`); err != nil {
		return nil, err
	}
	return calcServers, nil
}

func GetCalcServer(db *sqlx.DB, ip string) (CalcServers, error) {
	var calcServer CalcServers
	if err := db.Get(&calcServer, "SELECT * FROM calc_server_table WHERE ip=?", ip); err != nil {
		return calcServer, err
	}
	return calcServer, nil
}

func SetCalcServer(db *sqlx.DB, ip string, port string) error {
	_, err := db.NamedExec(`INSERT INTO calc_server_table (ip, port, status) VALUES (:ip,:port,:status)`,
		map[string]interface{}{
			"ip":     ip,
			"port":   port,
			"status": "active",
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCalcServerStatus(db *sqlx.DB, ip string, status string) error {
	_, err := db.NamedExec(`UPDATE calc_server_table SET status=:status WHERE ip=:ip`,
		map[string]interface{}{"ip": ip, "status": status},
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCalcServer(db *sqlx.DB, ip string) error {
	_, err := db.NamedExec(`DELETE FROM calc_server_table WHERE ip=:ip`, map[string]interface{}{"ip": ip})
	if err != nil {
		return err
	}
	return nil
}
