package main

import (
	"conda/modules"
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

	hub := modules.NewHub()
	go hub.Run()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Health(w, r)
	})
	r.Methods(http.MethodGet).Path("/connect").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.Connect(w, r, hub, db)
	})
	r.Methods(http.MethodGet).Path("/conda").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.EnvInfo(w, r)
	})
	r.Methods(http.MethodPut).Path("/upload").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controllers.FileUpload(w, r, db)
	})

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
	/*
		http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
				api.Connect(w, r, db)
			})

			http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
				api.UploadHandler(w, r, db)
			})
			http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
				api.KillHandler(w, r, db)
			})
			http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
				api.DeleteHandler(w, r, db)
			})
			http.HandleFunc("/trash", func(w http.ResponseWriter, r *http.Request) {
				api.TrashHandler(w, r, db)
			})
			http.HandleFunc("/env_info", func(w http.ResponseWriter, r *http.Request) {
				api.EnvInfoHandler(w, r)
			})
			http.HandleFunc("/host_status", func(w http.ResponseWriter, r *http.Request) {
				api.HostStatus(w, r)
			})
			http.HandleFunc("/process_status", func(w http.ResponseWriter, r *http.Request) {
				api.ProcessStatus(w, r, db)
			})
			http.HandleFunc("/programs/", func(w http.ResponseWriter, r *http.Request) {
				api.Explorer(w, r, db)
			})
			go api.ProcessStatusKernel()

			// サーバー
			http.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
				api.Servers(w, r, db)
			})
			http.HandleFunc("/exec_once", func(w http.ResponseWriter, r *http.Request) {
				api.ExecOnce(w, r, db)
			})

			log.Println("conda start")
			log.Fatal(http.ListenAndServe(":5984", nil))
	*/
}
