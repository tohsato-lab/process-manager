package modules

import (
	"database/sql"
	"log"
	"time"
)

//Process is struct
type Process struct {
	ID           string
	UseVram      float32
	Status       string
	Filename     string
	StartDate    *time.Time
	CompleteDate *time.Time
}

// RegistProcess データベースに新規登録
func RegistProcess(db *sql.DB, process *Process) {
	ins, err := db.Prepare("INSERT INTO process_table (id, use_vram, status, filename) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec("1qwdff8hf", "2", "ready", "test")
	if err != nil {
		log.Fatal(err)
	}
}

// UpdataProcess プロセスの更新
func UpdataProcess(db *sql.DB) {

}
