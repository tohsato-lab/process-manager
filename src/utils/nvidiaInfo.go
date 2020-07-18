package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// GetUsedVRAM 現在のVRAMの状況を取得
func GetUsedVRAM() float32 {
	cmdstr := "nvidia-smi | grep 'Default' | awk '{print $11-$9}'"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	result := strings.TrimRight(string(out), "\n")
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		fmt.Println("Comvert Error.")
	}
	return float32(value)
}

// GetTotalVRAM 現在のVRAMの状況を取得
func GetTotalVRAM() float32 {
	cmdstr := "nvidia-smi | grep 'Default' | awk '{print $11}' | sed 's@MiB@@g'"
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	result := strings.TrimRight(string(out), "\n")
	if err != nil {
		fmt.Println("Command Exec Error.")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		fmt.Println("Comvert Error.")
	}
	return float32(value)
}
