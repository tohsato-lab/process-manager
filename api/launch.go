package api

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Launch 起動スクリプト実行
func Launch(shellPath string) {
	os.Chdir(shellPath)
	cmd := exec.Command("bash", "launch.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	var buf bytes.Buffer
	multiWriter := io.MultiWriter(&buf, os.Stdout)
	cmd.Stdout = multiWriter

	fmt.Println(buf.String())
}
