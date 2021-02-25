package api

import (
	"fmt"
	"github.com/janberktold/sse"
	"net/http"
	"time"

	"process-manager-server/utils"
)

type Status struct {
	RAM  float64
	VRAM float64
}

// HostStatus ホストの情報を配信
func HostStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	conn, _ := sse.Upgrade(w, r)
	for {
		time.Sleep(1 * time.Second)
		totalVRAM := utils.GetTotalVRAM()
		totalRAM := utils.GetTotalRAM()
		var vram = 0.0
		var ram = 0.0
		if totalVRAM != 0 {
			vram = float64(utils.GetUsedVRAM() / totalVRAM)
		}
		if totalRAM != 0 {
			ram = float64(utils.GetUsedRAM() / totalRAM)
		}
		if err := conn.WriteJson(&Status{
			RAM:  ram,
			VRAM: vram,
		}); err != nil {
			fmt.Println(err)
			conn.Close()
			return
		}
	}
}
