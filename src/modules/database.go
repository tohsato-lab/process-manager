package modules

import (
	"database/sql"
	"fmt"
	"log"

	"../utils"
)

// RegistProcess データベースに新規登録
func RegistProcess(db *sql.DB, process utils.Process) {
	ins, err := db.Prepare("INSERT INTO process_table (id, use_vram, status, filename) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec(process.ID, process.UseVram, process.Status, process.Filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("registed")
	utils.BroadcastProcess <- process
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
		usedVRAM := 0
		err := db.QueryRow("SELECT SUM(use_vram) FROM process_table WHERE use_vram > 0 AND status = ?", "working").Scan(&usedVRAM)
		if err != nil {
			panic(err.Error())
		}
		// メモリに空きがある場合
		if vramTotal-float32(usedVRAM) >= 0 {
			var process utils.Process
			dbReady.Scan(&process.ID, &process.UseVram)
			statusUpdate, err := db.Prepare("UPDATE process_table SET status=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			defer statusUpdate.Close()

			if statusUpdate.Exec("working", process.ID); err != nil {
				panic(err.Error())
			}
			StartProcess(db, process.ID)
		} else {
			break
		}
	}

	// vramが0と設定されたプロセスを逐次実行
	var countReady int
	err = db.QueryRow("SELECT COUNT(*) FROM process_table WHERE use_vram = 0 AND status = ?", "ready").Scan(&countReady)
	if err != nil {
		panic(err.Error())
	}
	var countWorking int
	err = db.QueryRow("SELECT COUNT(*) FROM process_table WHERE status = ?", "working").Scan(&countWorking)
	if err != nil {
		panic(err.Error())
	}
	if countReady != 0 && countWorking == 0 {
		var process utils.Process
		err = db.QueryRow("SELECT id, use_vram FROM process_table WHERE use_vram = 0 AND status = ?", "ready").Scan(&process.ID, &process.UseVram)
		if err != nil {
			panic(err.Error())
		}

		statusUpdate, err := db.Prepare("UPDATE process_table SET status=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		defer statusUpdate.Close()

		if statusUpdate.Exec("working", process.ID); err != nil {
			panic(err.Error())
		}
		StartProcess(db, process.ID)
	}
}

// StartProcess プロセス実行
func StartProcess(db *sql.DB, id string) {
	go func() {
		Execute("../programs/" + id)
		ComplateProcess(db, id)
	}()
}

// ComplateProcess プロセス終了時にデータベースを更新
func ComplateProcess(db *sql.DB, id string) {
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer statusUpdate.Close()

	if statusUpdate.Exec("complate", id); err != nil {
		panic(err.Error())
	}
	UpdataAllProcess(db)
}
