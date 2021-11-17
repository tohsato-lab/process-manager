package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "docker:docker@tcp(mysql_host:3306)/process_manager_db?parseTime=true"
)

func main() {
	db, err := sql.Open(DriverName, DataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()



}
