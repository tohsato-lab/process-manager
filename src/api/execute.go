package api

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Execute 起動スクリプト実行
func Execute(shellPath string) {
	cmd := exec.Command("bash", "execute.sh", shellPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	var buf bytes.Buffer
	multiWriter := io.MultiWriter(&buf, os.Stdout)
	cmd.Stdout = multiWriter

	fmt.Println(buf.String())
}
