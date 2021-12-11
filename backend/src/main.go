package main

import (
	"backend/modules"
	"backend/repository"
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

func initDB(db *sqlx.DB) {
	servers, err := repository.GetActiveCalcServers(db)
	if err != nil {
		log.Println(err)
		return
	}
	for _, server := range servers {
		err := repository.UpdateCalcServerStatus(db, server.IP, "stop")
		if err != nil {
			log.Println(err)
			return
		}
	}
}

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

	initDB(db)
	modules.NewHub()

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
	r.Methods(http.MethodPut).Path("/process").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.EntryProcess(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/kill").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.KillProcess(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/trash").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.TrashAllProcess(w, r, db)
	})
	r.Methods(http.MethodPost).Path("/trash").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.TrashProcess(w, r, db)
	})
	r.Methods(http.MethodDelete).Path("/trash").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.DeleteProcess(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/connect").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Connect(w, r, db)
	})
	r.Methods(http.MethodGet).Path("/log").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.ProcessLog(w, r, db)
	})
	r.PathPrefix("/data/").Handler(controllers.SpaHandler{StaticPath: "../../"})

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

	log.Println("backend start")
	log.Fatal(http.ListenAndServe(":5983", c.Handler(r)))
}
