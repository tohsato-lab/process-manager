package modules

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"os/exec"
	"strconv"

	"conda/repository"
)

// Execute 起動スクリプト実行
func Execute(db *sqlx.DB, id string, targetFile string, envName string) string {
	//実行
	cmd := exec.Command("bash", "scripts/execute.sh", "../../data/"+id, targetFile, envName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "unknown:" + err.Error()
	}

	// PID登録
	err := repository.SetPID(db, id, cmd.Process.Pid)
	if err != nil {
		log.Println(err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	var status string
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
