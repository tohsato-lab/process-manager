package modules

import (
	"backend/repository"
	"github.com/jmoiron/sqlx"
	"os"
	"os/exec"
)

func DownloadLogs(db *sqlx.DB, processID string) (int, error) {
	process, err := repository.GetProcess(db, processID)
	if err != nil {
		return 1, err
	}
	cmd := exec.Command("bash", "scripts/lftp.sh", process.ServerIP, processID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return 1, err
	}
	if err := cmd.Wait(); err != nil {
		return 1, err
	}
	signal := cmd.ProcessState.ExitCode()
	return signal, nil
}
