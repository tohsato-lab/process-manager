package modules

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"process-manager-server/utils"
)

// GetAllProcess プロセス一覧取得
func GetAllProcess(db *sql.DB) []utils.Process {
	fmt.Println("### GetAllProcess")
	var processes []utils.Process

	dbSelect, err := db.Query("SELECT id, use_vram, status, filename, start_date, complete_date, exec_count, comment FROM process_table ORDER BY start_date DESC")
	if err != nil {
		fmt.Println(err)
	}
	defer dbSelect.Close()
	for dbSelect.Next() {
		var process utils.Process
		var startDate sql.NullTime
		var completeDate sql.NullTime
		var execCount sql.NullInt32
		var comment sql.NullString
		if err := dbSelect.Scan(&process.ID, &process.UseVram, &process.Status, &process.Filename, &startDate, &completeDate, &execCount, &comment); err != nil {
			fmt.Println(err)
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if startDate.Valid {
			process.StartDate = startDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if completeDate.Valid {
			process.CompleteDate = completeDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if execCount.Valid {
			process.ExecCount = execCount.Int32
		}
		if comment.Valid {
			process.Comment = comment.String
		}
		processes = append(processes, process)
	}
	return processes
}

// RegisterProcess データベースに新規登録
func RegisterProcess(db *sql.DB, process utils.Process) {
	fmt.Println("### RegisterProcess")
	ins, err := db.Prepare("INSERT INTO process_table (id, use_vram, status, filename, targetfile, env_name, exec_count, comment) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ins.Exec(process.ID, process.UseVram, process.Status, process.Filename, process.TargetFile, process.EnvName, process.ExecCount, process.Comment); err != nil {
		fmt.Println(err)
	}
	utils.BroadcastProcess <- GetAllProcess(db)
	UpdateAllProcess(db)
}

// UpdateAllProcess プロセスの更新
func UpdateAllProcess(db *sql.DB) {
	fmt.Println("### UpdateAllProcess")
	vramTotal := utils.GetTotalVRAM()

	// 0より大きく設定されたプロセスを先に実行
	dbReady, err := db.Query("SELECT id, use_vram, targetfile, env_name FROM process_table WHERE use_vram > 0 AND exec_count > 0 AND status = ? ORDER BY use_vram", "ready")
	if err != nil {
		fmt.Println(err)
	}
	defer dbReady.Close()

	for dbReady.Next() {
		usedVRAM := 0
		err := db.QueryRow("SELECT IFNULL(SUM(use_vram), 0) FROM process_table WHERE use_vram > 0 AND exec_count > 0 AND status = ?", "running").Scan(&usedVRAM)
		if err != nil {
			fmt.Println(err)
		}
		var process utils.Process
		if err := dbReady.Scan(&process.ID, &process.UseVram, &process.TargetFile, &process.EnvName, &process.ExecCount); err != nil {
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
	if err := db.QueryRow("SELECT COUNT(*) FROM process_table WHERE use_vram = 0 AND exec_count > 0 AND status = ?", "ready").Scan(&countReady); err != nil {
		fmt.Println(err)
	}
	var countWorking int
	if err := db.QueryRow("SELECT COUNT(*) FROM process_table WHERE status = ?", "running").Scan(&countWorking); err != nil {
		fmt.Println(err)
	}
	if countReady != 0 && countWorking == 0 {
		var process utils.Process
		if err := db.QueryRow(
			"SELECT id, targetfile, env_name FROM process_table WHERE use_vram = 0 AND exec_count > 0 AND status = ?", "ready",
		).Scan(&process.ID, &process.TargetFile, &process.EnvName); err != nil {
			fmt.Println(err)
		}
		StartProcess(db, process.ID, process.TargetFile, process.EnvName)
	}
	utils.BroadcastProcess <- GetAllProcess(db)
}

// RegisterPID データベースにPID登録
func RegisterPID(db *sql.DB, id string, pid int) {
	fmt.Println("### RegisterPID")

	statusUpdate, err := db.Prepare("UPDATE process_table SET pid=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := statusUpdate.Exec(strconv.Itoa(pid), id); err != nil {
		fmt.Println(err)
	}
	if err := statusUpdate.Close(); err != nil {
		fmt.Println(err)
	}
}

// StartProcess プロセス実行
func StartProcess(db *sql.DB, id string, targetFile string, envName string) {
	fmt.Println("### StartProcess")

	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, start_date=?, complete_date=NULL WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := statusUpdate.Exec("running", time.Now(), id); err != nil {
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
	UpdateAllProcess(db)
}

// CompleteProcess プロセス終了時にデータベースを更新
func CompleteProcess(db *sql.DB, id string, status string) {
	fmt.Println("### CompleteProcess")

	var execCount int
	if err := db.QueryRow("SELECT exec_count FROM process_table WHERE id = ?", id).Scan(&execCount); err != nil {
		fmt.Println(err)
	}
	statusUpdate, err := db.Prepare("UPDATE process_table SET status=?, complete_date=?, exec_count=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	if execCount-1 > 0 {
		status = "ready"
	}
	if _, err := statusUpdate.Exec(status, time.Now(), execCount-1, id); err != nil {
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
	fmt.Println("### DeleteProcess")

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
	UpdateAllProcess(db)
}
