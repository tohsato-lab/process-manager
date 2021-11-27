package modules

import (
	"backend/utils"
	"database/sql"
	"fmt"
)

func GetServers(db *sql.DB) []utils.Servers {
	dbSelect, err := db.Query(`SELECT ip, port, status FROM servers`)
	if err != nil {
		fmt.Println(err)
	}
	var servers []utils.Servers
	defer dbSelect.Close()
	for dbSelect.Next() {
		var server utils.Servers
		if err := dbSelect.Scan(
			&server.IP, &server.Port, &server.Status,
		); err != nil {
			fmt.Println(err)
		}
		servers = append(servers, server)
	}
	return servers
}
