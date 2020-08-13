package modules

import (
	"database/sql"
	"os"
	"os/exec"
	"strconv"
)

// Execute 起動スクリプト実行
func Execute(db *sql.DB, id string, targetfile string, envName string) string {
	//実行
	cmd := exec.Command("bash", "execute.sh", "../data/programs/"+id, targetfile, envName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	// PID登録
	statusUpdate, err := db.Prepare("UPDATE process_table SET pid=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer statusUpdate.Close()
	if statusUpdate.Exec(strconv.Itoa(cmd.Process.Pid), id); err != nil {
		panic(err.Error())
	}

	cmd.Wait()

	status := ""
	signal := cmd.ProcessState.ExitCode()
	switch signal {
	case 0:
		status = "complete"
	case 1:
		status = "error"
	case 143:
		status = "killed"
	default:
		status = "unknown:" + strconv.Itoa(signal)
	}
	return status
}
