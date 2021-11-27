package modules

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"conda/utils"
)

// GetProcesses プロセス一覧取得
func GetProcesses(db *sql.DB) []utils.Process {
	fmt.Println("### GetProcesses")

	var processes []utils.Process

	dbSelect, err := db.Query(
		`SELECT id, process_name, status, env_name, target_file, start_date, complete_date, comment 
			   FROM main_processes WHERE !in_trash ORDER BY upload_date DESC`,
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
			&process.ID,
			&process.ProcessName,
			&process.Status,
			&process.EnvName,
			&process.TargetFile,
			&startDate,
			&completeDate,
			&comment,
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
		`INSERT INTO main_processes
               (id, process_name, status, target_file, env_name, comment, upload_date, in_trash, is_home, server_ip) 
               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ins.Exec(
		process.ID,
		process.ProcessName,
		process.Status,
		process.TargetFile,
		process.EnvName,
		process.Comment,
		time.Now(),
		process.InTrash,
		process.IsHome,
		process.ServerIP,
	); err != nil {
		fmt.Println(err)
	}
	UpdateAllProcess(db)
}

// UpdateAllProcess プロセスの更新
func UpdateAllProcess(db *sql.DB) {
	fmt.Println("### UpdateAllProcess")

	// is_home process
	fmt.Println("check home process")
	var countRunning int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM main_processes WHERE status = 'running' AND is_home = true`,
	).Scan(&countRunning); err != nil {
		fmt.Println(err)
	}
	if countRunning == 0 {
		dbSelect, err := db.Query(
			`SELECT id, target_file, env_name 
				   FROM main_processes 
				   WHERE status = 'ready' AND in_trash = false AND is_home = true 
				   ORDER BY upload_date LIMIT 1`,
		)
		if err != nil {
			fmt.Println(err)
		}
		defer dbSelect.Close()
		for dbSelect.Next() {
			var process utils.Process
			if err := dbSelect.Scan(&process.ID, &process.TargetFile, &process.EnvName); err != nil {
				fmt.Println(err)
			}
			StartProcess(db, process.ID, process.TargetFile, process.EnvName)
		}
	}

	// 外部からのプロセス
	fmt.Println("check not home process")
	dbSelect, err := db.Query(
		`SELECT id, target_file, env_name 
			   FROM main_processes 
			   WHERE status = 'ready' AND in_trash = false AND is_home = false 
			   ORDER BY upload_date`,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer dbSelect.Close()
	for dbSelect.Next() {
		var process utils.Process
		if err := dbSelect.Scan(&process.ID, &process.TargetFile, &process.EnvName); err != nil {
			fmt.Println(err)
		}
		StartProcess(db, process.ID, process.TargetFile, process.EnvName)
	}

	utils.BroadcastProcesses <- GetProcesses(db)
}

// RegisterPID データベースにPID登録
func RegisterPID(db *sql.DB, id string, pid int) {
	fmt.Println("### RegisterPID")

	statusUpdate, err := db.Prepare(`UPDATE main_processes SET pid=? WHERE id=?`)
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
		`UPDATE main_processes SET status=?, start_date=?, complete_date=NULL WHERE id=?`,
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

	statusUpdate, err := db.Prepare(`UPDATE main_processes SET status=?, complete_date=? WHERE id=?`)
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

	dbDelete, err := db.Prepare(`DELETE FROM main_processes WHERE id = ?`)
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

// GetTrashProcesses ゴミ箱プロセス一覧取得
func GetTrashProcesses(db *sql.DB) []utils.Process {
	fmt.Println("### GetTrashProcesses")

	var processes []utils.Process

	dbSelect, err := db.Query(
		`SELECT id, status, env_name, target_file, start_date, complete_date, comment 
			   FROM main_processes WHERE in_trash 
			   ORDER BY upload_date DESC`,
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
			&process.ID, &process.Status, &process.EnvName, &process.TargetFile, &startDate, &completeDate, &comment,
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

// TrashProcess ゴミ箱
func TrashProcess(db *sql.DB, id string) {
	fmt.Println("### TrashProcess")

	var trashStatus bool

	if err := db.QueryRow(
		`SELECT in_trash FROM main_processes WHERE id=?`, id,
	).Scan(&trashStatus); err != nil {
		fmt.Println(err)
	}
	dbDelete, err := db.Prepare(`UPDATE main_processes SET in_trash=? WHERE id=?`)
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

// GetServers サーバーリスト取得
func GetServers(db *sql.DB) []utils.Servers {
	dbSelect, err := db.Query(`SELECT ip, port, status FROM servers`)
	if err != nil {
		fmt.Println(err)
	}
	var servers []utils.Servers
	defer dbSelect.Close()
	for dbSelect.Next() {
		var server utils.Servers
		if err := dbSelect.Scan(
			&server.IP, &server.Port, &server.Status,
		); err != nil {
			fmt.Println(err)
		}
		servers = append(servers, server)
	}
	return servers
}

// RegisterServer サーバー登録
func RegisterServer(db *sql.DB, ip string, port string) {
	ins, err := db.Prepare(`INSERT INTO servers (ip, port, status) VALUES (?, ?, ?)`)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := ins.Exec(ip, port, "arrive"); err != nil {
		fmt.Println(err)
	}
}

func DeleteServer(db *sql.DB, ip string) {
	dbDelete, err := db.Prepare(`DELETE FROM servers WHERE ip = ?`)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := dbDelete.Exec(ip); err != nil {
		fmt.Println(err)
	}
	if err := dbDelete.Close(); err != nil {
		fmt.Println(err)
	}
}
