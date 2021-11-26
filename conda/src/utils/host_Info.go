package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetUsedRAM 現在のRAMの状況を取得
func GetUsedRAM() float32 {
	cmdStr := "free | grep Mem | awk '{print $3}'"
	out, err := exec.Command("sh", "-c", cmdStr).Output()
	result := strings.TrimRight(string(out), "\n")
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		fmt.Println("Convert Error.")
	}
	return float32(value)
}

// GetTotalRAM RAMの容量を取得
func GetTotalRAM() float32 {
	cmdStr := "free | grep Mem | awk '{print $2}'"
	out, err := exec.Command("sh", "-c", cmdStr).Output()
	result := strings.TrimRight(string(out), "\n")
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		fmt.Println("Convert Error.")
	}
	return float32(value)
}
