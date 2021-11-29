package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"

	"backend/utils"
)

func GetCalcServers(db *sqlx.DB) ([]byte, error) {
	var calcServers []utils.CalcServers
	if err := db.Select(&calcServers, `SELECT * FROM servers`); err != nil {
		return nil, err
	}
	contents, err := json.Marshal(calcServers)
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func SetCalcServer(db *sqlx.DB, ip string, port string) error {
	_, err := db.NamedExec(`INSERT INTO servers (ip, port, status) VALUES (:ip,:port,:status)`,
		map[string]interface{}{
			"ip":     ip,
			"port":   port,
			"status": "arrive",
		},
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
