package modules

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Execute 起動スクリプト実行
func Execute(db *sql.DB, id string, targetFile string, envName string) string {
	//実行
	cmd := exec.Command("bash", "scripts/execute.sh", "../../data/programs/"+id, targetFile, envName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "unknown:" + err.Error()
	}

	// PID登録
	RegisterPID(db, id, cmd.Process.Pid)

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	status := ""
	signal := cmd.ProcessState.ExitCode()
	switch signal {
	case 0:
		status = "complete"
	case 1:
		status = "error"
	case 128 + 15:
		status = "killed"
	default:
		status = "unknown:" + strconv.Itoa(signal)
	}
	return status
}
