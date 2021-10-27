package api

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// ServerRegister sshでマウント
func ServerRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if r.Method != "GET" {
		_, _ = fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	ip, err := getIP(r)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("No valid ip"))
	}
	fmt.Println(ip)

	mountPoint := "../../data/" + ip
	if err := os.MkdirAll(mountPoint, 0777); err != nil {
		_, _ = fmt.Fprintln(w, "ディレクトリ生成に失敗しました。"+err.Error())
		return
	}

	command := "echo docker | sshfs -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no docker@" +
		ip + ":process-manager/data " + mountPoint + " -o workaround=rename -o password_stdin -o follow_symlinks -p 8022"
	if _, err := exec.Command("sh", "-c", command).Output(); err != nil {
		_, _ = fmt.Fprintln(w, "マウントに失敗しました。"+err.Error())
		return
	}

}
func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", err
}
