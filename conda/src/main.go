package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"log"
	"net/http"

	"conda/controllers"
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
	r.Methods(http.MethodGet).Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Health(w, r)
	})
	r.Methods(http.MethodGet).Path("/connect").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Connect(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/conda").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.EnvInfo(w, r)
	})
	r.Methods(http.MethodPut).Path("/upload").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.FileUpload(w, r, db)
	})
	r.Methods(http.MethodDelete).Path("/delete").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteFile(w, r, db)
	})
	r.PathPrefix("/log/").Handler(controllers.SpaHandler{StaticPath: "../../"})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodOptions,
		},
	})

	log.Println("conda start")
	log.Fatal(http.ListenAndServe(":5984", c.Handler(r)))
}
