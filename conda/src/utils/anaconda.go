package utils

import (
	"log"
	"os/exec"
	"strings"
)

// GetCondaEnv Anacondaのenv一覧を取得
func GetCondaEnv() ([]string, error) {
	cmdStr := "conda info -e | grep 'conda' | grep -v '#' | awk '{print $1}'"
	out, err := exec.Command("sh", "-c", cmdStr).Output()
	if err != nil {
		log.Println("Command Exec Error.")
		return nil, err
	}
	envArray := strings.Split(string(out), "\n")
	return envArray[:len(envArray)-1], nil
}
