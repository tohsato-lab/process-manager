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

func execute(db *sqlx.DB, id string, targetFile string, args string, envName string) (string, error) {
	//実行
	cmd := exec.Command("bash", "scripts/execute.sh", "../../log/"+id, targetFile, args, envName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "unknown:" + err.Error(), err
	}

	// PID登録
	err := repository.SetPID(db, id, cmd.Process.Pid)
	if err != nil {
		return "error", err
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
	return status, nil
}

func killCMD(db *sqlx.DB, processID string) (string, error) {
	process, err := repository.GetProcess(db, processID)
	if err != nil {
		return "", err
	}
	if process.PID.Valid {
		cmd := "kill `ps ho pid --ppid=" + strconv.Itoa(int(process.PID.Int32)) + "`"
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return "", err
		}
		return "killed", nil
	}
	return process.Status, nil
}

func DeleteCMD(processID string) error {
	log.Println("delete")
	cmd := "rm -rf ../../log/" + processID + "/"
	log.Println(cmd)
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return err
	}
	return nil
}
