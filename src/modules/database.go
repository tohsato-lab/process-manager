package modules

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"../utils"
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

// UpdataAllProcess プロセスの更新
func UpdataAllProcess(db *sql.DB) {
	vramTotal := utils.GetTotalVRAM()
	fmt.Println(vramTotal)

	// 0より大きく設定されたプロセスを先に実行
	dbReady, err := db.Query("SELECT id, use_vram FROM process_table WHERE use_vram > 0 AND status = ? ORDER BY use_vram", "ready")
	if err != nil {
		panic(err.Error())
	}
	defer dbReady.Close()

	for dbReady.Next() {
		dbUsedVRAM, err := db.Query("SELECT SUM(use_vram) FROM process_table WHERE use_vram > 0 AND status = ?", "working")
		if err != nil {
			panic(err.Error())
		}
		defer dbUsedVRAM.Close()
		dbUsedVRAM.Next()
		usedVRAM := 0
		dbUsedVRAM.Scan(&usedVRAM)
		// メモリに空きがある場合
		if vramTotal-float32(usedVRAM) >= 0 {
			var process Process
			dbReady.Scan(&process.ID, &process.UseVram)
			statusUpdate, err := db.Prepare("UPDATE process_table SET status=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			defer statusUpdate.Close()

			if statusUpdate.Exec("working", process.ID); err != nil {
				panic(err.Error())
			}
			Execute("../programs/" + process.ID)
		} else {
			break
		}
	}

}
