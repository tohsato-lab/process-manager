package modules

import (
	"database/sql"
	"fmt"
	"time"

	"process-manager-server/utils"
)

// GetAllProcess プロセス一覧取得
func GetAllProcess(db *sql.DB) []utils.Process {
	var processes []utils.Process

	dbSelect, err := db.Query("SELECT id, use_vram, status, filename, start_date, complete_date FROM process_table ORDER BY start_date DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer dbSelect.Close()
	for dbSelect.Next() {
		var process utils.Process
		var startDate sql.NullTime
		var completeDate sql.NullTime
		if err := dbSelect.Scan(&process.ID, &process.UseVram, &process.Status, &process.Filename, &startDate, &completeDate); err != nil {
			fmt.Println(err)
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if startDate.Valid {
			process.StartDate = startDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if completeDate.Valid {
			process.CompleteDate = completeDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		processes = append(processes, process)
	}
	return processes
}

// RegisterProcess データベースに新規登録
func RegisterProcess(db *sql.DB, process utils.Process) {
	ins, err := db.Prepare("INSERT INTO process_table (id, use_vram, status, filename, targetfile, env_name) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ins.Exec(process.ID, process.UseVram, process.Status, process.Filename, process.TargetFile, process.EnvName); err != nil {
		fmt.Println(err)
	}
	utils.BroadcastProcess <- GetAllProcess(db)
}

// UpdateAllProcess プロセスの更新
func UpdateAllProcess(db *sql.DB) {

	vramTotal := utils.GetTotalVRAM()

	// 0より大きく設定されたプロセスを先に実行
	dbReady, err := db.Query("SELECT id, use_vram, targetfile, env_name FROM process_table WHERE use_vram > 0 AND status = ? ORDER BY use_vram", "ready")
	if err != nil {
		fmt.Println(err)
	}
	defer dbReady.Close()

	for dbReady.Next() {
		usedVRAM := 0
		err := db.QueryRow("SELECT IFNULL(SUM(use_vram), 0) FROM process_table WHERE use_vram > 0 AND status = ?", "working").Scan(&usedVRAM)
		if err != nil {
			fmt.Println(err)
		}
		var process utils.Process
		if err := dbReady.Scan(&process.ID, &process.UseVram, &process.TargetFile, &process.EnvName); err != nil {
			fmt.Println(err)
		}

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
	if err := db.QueryRow("SELECT COUNT(*) FROM process_table WHERE use_vram = 0 AND status = ?", "ready").Scan(&countReady); err != nil {
		fmt.Println(err)
	}
	var countWorking int
	if err := db.QueryRow("SELECT COUNT(*) FROM process_table WHERE status = ?", "working").Scan(&countWorking); err != nil {
		fmt.Println(err)
	}
	if countReady != 0 && countWorking == 0 {
		var process utils.Process
		if err := db.QueryRow(
			"SELECT id, targetfile, env_name FROM process_table WHERE use_vram = 0 AND status = ?", "ready",
		).Scan(&process.ID, &process.TargetFile, &process.EnvName); err != nil {
			fmt.Println(err)
		}
		StartProcess(db, process.ID, process.TargetFile, process.EnvName)
	}
}

// StartProcess プロセス実行
func StartProcess(db *sql.DB, id string, targetFile string, envName string) {
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, start_date=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := statusUpdate.Exec("working", time.Now(), id); err != nil {
		fmt.Println(err)
	}

	if err := statusUpdate.Close(); err != nil {
		fmt.Println(err)
	}

	go func() {
		status := Execute(db, id, targetFile, envName)
		CompleteProcess(db, id, status)
	}()

	utils.BroadcastProcess <- GetAllProcess(db)
}

// CompleteProcess プロセス終了時にデータベースを更新
func CompleteProcess(db *sql.DB, id string, status string) {
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, complete_date=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := statusUpdate.Exec(status, time.Now(), id); err != nil {
		fmt.Println(err)
	}

	if err := statusUpdate.Close(); err != nil {
		fmt.Println(err)
	}

	utils.BroadcastProcess <- GetAllProcess(db)

	UpdateAllProcess(db)
}

// DeleteProcess リストからプロセスを削除
func DeleteProcess(db *sql.DB, id string) {
	dbDelete, err := db.Prepare("DELETE FROM process_table WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}

	if _, err := dbDelete.Exec(id); err != nil {
		fmt.Println(err)
	}

	if err := dbDelete.Close(); err != nil {
		fmt.Println(err)
	}

	utils.BroadcastProcess <- GetAllProcess(db)
}
