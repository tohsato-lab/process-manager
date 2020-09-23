package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetCondaEnv Anacondaのenv一覧を取得
func GetCondaEnv() []string {
	cmdstr := "conda info -e | grep anaconda3 | awk '{print $1}'"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	return strings.Split(string(out), "\n")[:4]
}
