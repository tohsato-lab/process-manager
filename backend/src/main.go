package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"log"
	"net/http"

	"backend/controllers"
)

const (
	// DriverName ドライバ名(mysql固定)
	DriverName = "mysql"
	// DataSourceName user:password@tcp(container-name:port)/dbname
	DataSourceName = "docker:docker@tcp(mysql_host:3306)/process_manager_db?parseTime=true"
)

func main() {
	db, err := sqlx.Open(DriverName, DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/calculator").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ServerStatus(w, r, db)
	})
	r.Methods(http.MethodPost).Path("/calculator").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.JoinServer(w, r, db)
	})
	r.Methods(http.MethodDelete).Path("/calculator").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteServer(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/conda").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Envs(w, r, db)
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	})

	log.Println("backend start")
	log.Fatal(http.ListenAndServe(":5983", c.Handler(r)))
}
