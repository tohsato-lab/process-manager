package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetUsedVRAM 現在のVRAMの状況を取得
func GetUsedVRAM() float32 {
	cmdStr := "nvidia-smi | grep 'Default' | awk '{print $9}' | sed 's@MiB@@g'"
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

// GetTotalVRAM VRAMの容量を取得
func GetTotalVRAM() float32 {
	cmdStr := "nvidia-smi | grep 'Default' | awk '{print $11}' | sed 's@MiB@@g'"
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
