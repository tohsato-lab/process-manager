package modules

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"process-manager-server/utils"
)

// GetProcesses プロセス一覧取得
func GetProcesses(db *sql.DB) []utils.Process {
	fmt.Println("### GetProcesses")

	var processes []utils.Process

	dbSelect, err := db.Query(
		"SELECT id, status, filename, env_name, target_file, start_date, complete_date, comment FROM main_processes WHERE !in_trash ORDER BY upload_date DESC",
	)
	if err != nil {
		fmt.Println(err)
	}
	defer dbSelect.Close()
	for dbSelect.Next() {
		var process utils.Process
		var startDate sql.NullTime
		var completeDate sql.NullTime
		var comment sql.NullString
		if err := dbSelect.Scan(
			&process.ID, &process.Status, &process.Filename, &process.EnvName, &process.TargetFile, &startDate, &completeDate, &comment,
		); err != nil {
			fmt.Println(err)
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if startDate.Valid {
			process.StartDate = startDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
		}
		if completeDate.Valid {
			process.CompleteDate = completeDate.Time.In(jst).Format("2006年01月02日 15時04分05秒")
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
	ins, err := db.Prepare(
		"INSERT INTO main_processes (id, status, filename, target_file, env_name, comment, upload_date, in_trash) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ins.Exec(
		process.ID,
		process.Status,
		process.Filename,
		process.TargetFile,
		process.EnvName,
		process.Comment,
		time.Now(),
		process.InTrash,
	); err != nil {
		fmt.Println(err)
	}
	UpdateAllProcess(db)
}

// UpdateAllProcess プロセスの更新
func UpdateAllProcess(db *sql.DB) {
	fmt.Println("### UpdateAllProcess")

	var countReady int
	var countRunning int

	if err := db.QueryRow(
		"SELECT COUNT(*) FROM main_processes WHERE status = ?", "ready",
	).Scan(&countReady); err != nil {
		fmt.Println(err)
	}
	if err := db.QueryRow(
		"SELECT COUNT(*) FROM main_processes WHERE status = ?", "running",
	).Scan(&countRunning); err != nil {
		fmt.Println(err)
	}

	if countReady != 0 && countRunning == 0 {
		var process utils.Process
		if err := db.QueryRow(
			"SELECT id, target_file, env_name FROM main_processes WHERE status = ? ORDER BY upload_date", "ready",
		).Scan(&process.ID, &process.TargetFile, &process.EnvName); err != nil {
			fmt.Println(err)
		}
		StartProcess(db, process.ID, process.TargetFile, process.EnvName)
	}
	utils.BroadcastProcesses <- GetProcesses(db)
}

// RegisterPID データベースにPID登録
func RegisterPID(db *sql.DB, id string, pid int) {
	fmt.Println("### RegisterPID")

	statusUpdate, err := db.Prepare("UPDATE main_processes SET pid=? WHERE id=?")
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

	statusUpdate, err := db.Prepare(
		"UPDATE main_processes SET status=?, start_date=?, complete_date=NULL WHERE id=?",
	)
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
		// TODO: OS判定？
		status := Execute(db, id, targetFile, envName)
		CompleteProcess(db, id, status)
	}()
	UpdateAllProcess(db)
}

// CompleteProcess プロセス終了時にデータベースを更新
func CompleteProcess(db *sql.DB, id string, status string) {
	fmt.Println("### CompleteProcess")

	statusUpdate, err := db.Prepare("UPDATE main_processes SET status=?, complete_date=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := statusUpdate.Exec(status, time.Now(), id); err != nil {
		fmt.Println(err)
	}
	if err := statusUpdate.Close(); err != nil {
		fmt.Println(err)
	}
	UpdateAllProcess(db)
}

// DeleteProcess リストからプロセスを削除
func DeleteProcess(db *sql.DB, id string) {
	fmt.Println("### DeleteProcess")

	dbDelete, err := db.Prepare("DELETE FROM main_processes WHERE id = ?")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := dbDelete.Exec(id); err != nil {
		fmt.Println(err)
	}
	if err := dbDelete.Close(); err != nil {
		fmt.Println(err)
	}
	UpdateAllProcess(db)
}

// TrashProcess ゴミ箱
func TrashProcess(db *sql.DB, id string) {
	fmt.Println("### TrashProcess")

	var trashStatus bool

	if err := db.QueryRow(
		"SELECT in_trash FROM main_processes WHERE id=?", id,
	).Scan(&trashStatus); err != nil {
		fmt.Println(err)
	}
	dbDelete, err := db.Prepare("UPDATE main_processes SET in_trash=? WHERE id=?")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := dbDelete.Exec(!trashStatus, id); err != nil {
		fmt.Println(err)
	}
	if err := dbDelete.Close(); err != nil {
		fmt.Println(err)
	}
	UpdateAllProcess(db)
}
