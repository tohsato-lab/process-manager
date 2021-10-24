package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetCondaEnv Anacondaのenv一覧を取得
func GetCondaEnv() []string {
	cmdStr := "conda info -e | grep 'conda' | grep -v '#' | awk '{print $1}'"
	out, err := exec.Command("sh", "-c", cmdStr).Output()
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	envArray := strings.Split(string(out), "\n")
	return envArray[:len(envArray)-1]
}
