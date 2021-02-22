package modules

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"../utils"
)

// GetAllProcess プロセス一覧取得
func GetAllProcess(db *sql.DB) []utils.Process {
	var processes []utils.Process

	dbSelect, err := db.Query("SELECT id, use_vram, status, filename, start_date, complete_date FROM process_table ORDER BY start_date DESC")
	if err != nil {
		panic(err.Error())
	}
	defer dbSelect.Close()
	for dbSelect.Next() {
		var process utils.Process
		var startDate time.Time
		var completeDate time.Time
		dbSelect.Scan(&process.ID, &process.UseVram, &process.Status, &process.Filename, &startDate, &completeDate)
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if !startDate.IsZero() {
			process.StartDate = startDate.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if !completeDate.IsZero() {
			process.CompleteDate = completeDate.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		processes = append(processes, process)
	}
	return processes
}

// RegistProcess データベースに新規登録
func RegistProcess(db *sql.DB, process utils.Process) {
	ins, err := db.Prepare("INSERT INTO process_table (id, use_vram, status, filename, targetfile, env_name) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	ins.Exec(process.ID, process.UseVram, process.Status, process.Filename, process.TargetFile, process.EnvName)
	if err != nil {
		log.Fatal(err)
	}
	utils.BroadcastProcess <- GetAllProcess(db)

}

// UpdataAllProcess プロセスの更新
func UpdataAllProcess(db *sql.DB) {

	vramTotal := utils.GetTotalVRAM()

	// 0より大きく設定されたプロセスを先に実行
	dbReady, err := db.Query("SELECT id, use_vram, targetfile, env_name FROM process_table WHERE use_vram > 0 AND status = ? ORDER BY use_vram", "ready")
	if err != nil {
		panic(err.Error())
	}
	defer dbReady.Close()

	for dbReady.Next() {
		usedVRAM := 0
		err := db.QueryRow("SELECT IFNULL(SUM(use_vram), 0) FROM process_table WHERE use_vram > 0 AND status = ?", "working").Scan(&usedVRAM)
		if err != nil {
			panic(err.Error())
		}
		var process utils.Process
		dbReady.Scan(&process.ID, &process.UseVram, &process.TargetFile, &process.EnvName)

		// メモリに空きがある場合
		fmt.Println(vramTotal - (float32(usedVRAM) + process.UseVram))
		if vramTotal-(float32(usedVRAM)+process.UseVram) >= 0 {
			StartProcess(db, process.ID, process.TargetFile, process.EnvName)
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
		err = db.QueryRow("SELECT id, targetfile, env_name FROM process_table WHERE use_vram = 0 AND status = ?", "ready").Scan(&process.ID, &process.TargetFile, &process.EnvName)
		if err != nil {
			panic(err.Error())
		}
		StartProcess(db, process.ID, process.TargetFile, process.EnvName)
	}
}

// StartProcess プロセス実行
func StartProcess(db *sql.DB, id string, targetfile string, envName string) {
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, start_date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer statusUpdate.Close()

	if _, err := statusUpdate.Exec("working", time.Now(), id); err != nil {
		panic(err.Error())
	}

	go func() {
		status := Execute(db, id, targetfile, envName)
		ComplateProcess(db, id, status)
	}()

	utils.BroadcastProcess <- GetAllProcess(db)
}

// ComplateProcess プロセス終了時にデータベースを更新
func ComplateProcess(db *sql.DB, id string, status string) {
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, complete_date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer statusUpdate.Close()

	if _, err := statusUpdate.Exec(status, time.Now(), id); err != nil {
		panic(err.Error())
	}

	utils.BroadcastProcess <- GetAllProcess(db)

	UpdataAllProcess(db)
}

// DeleteProcess リストからプロセスを削除
func DeleteProcess(db *sql.DB, id string) {
	dbDelete, err := db.Prepare("DELETE FROM process_table WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer dbDelete.Close()

	if _, err := dbDelete.Exec(id); err != nil {
		panic(err.Error())
	}

	utils.BroadcastProcess <- GetAllProcess(db)
}
