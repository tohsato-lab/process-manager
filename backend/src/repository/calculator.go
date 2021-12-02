package repository

import (
	"github.com/jmoiron/sqlx"

	"backend/utils"
)

func GetAllCalcServers(db *sqlx.DB) ([]utils.CalcServers, error) {
	var calcServers []utils.CalcServers
	if err := db.Select(&calcServers, `SELECT * FROM servers`); err != nil {
		return nil, err
	}
	return calcServers, nil
}

func GetActiveCalcServers(db *sqlx.DB) ([]utils.CalcServers, error) {
	var calcServers []utils.CalcServers
	if err := db.Select(&calcServers, `SELECT * FROM servers WHERE status='active'`); err != nil {
		return nil, err
	}
	return calcServers, nil
}

func GetCalcServerPort(db *sqlx.DB, ip string) (utils.CalcServers, error) {
	var calcServer utils.CalcServers
	if err := db.Get(&calcServer, "SELECT * FROM servers WHERE ip=?", ip); err != nil{
		return calcServer, err
	}
	return calcServer, nil
}

func SetCalcServer(db *sqlx.DB, ip string, port string) error {
	_, err := db.NamedExec(`INSERT INTO servers (ip, port, status) VALUES (:ip,:port,:status)`,
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
	_, err := db.NamedExec(`UPDATE servers SET status=:status WHERE ip=:ip`,
		map[string]interface{}{"ip": ip, "status": status},
	)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCalcServer(db *sqlx.DB, ip string) error {
	_, err := db.NamedExec(`DELETE FROM servers WHERE ip=:ip`, map[string]interface{}{"ip": ip})
	if err != nil {
		return err
	}
	return nil
}
