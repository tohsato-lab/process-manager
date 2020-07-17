package modules

import (
	"database/sql"
	"fmt"
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
	ins.Exec(process.ID, process.UseVram, process.Status, process.Filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("registed")
}

// UpdataProcess プロセスの更新
func UpdataProcess(db *sql.DB) {

}
