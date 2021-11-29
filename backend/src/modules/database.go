package modules

import (
	"backend/utils"
	"database/sql"
	"fmt"
)

func GetServers(db *sql.DB) []utils.CalcServers {
	dbSelect, err := db.Query(`SELECT ip, port, status FROM servers`)
	if err != nil {
		fmt.Println(err)
	}
	var servers []utils.CalcServers
	defer dbSelect.Close()
	for dbSelect.Next() {
		var server utils.CalcServers
		if err := dbSelect.Scan(
			&server.IP, &server.Port, &server.Status,
		); err != nil {
			fmt.Println(err)
		}
		servers = append(servers, server)
	}
	return servers
}
