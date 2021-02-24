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

// HostStatus GPUの情報を配信
func HostStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	conn, _ := sse.Upgrade(w, r)
	for {
		time.Sleep(1 * time.Second)
		totalVRAM := utils.GetTotalVRAM()
		var vram = 0.0
		if totalVRAM != 0 {
			vram = float64(utils.GetUsedVRAM() / totalVRAM)
		}
		if err := conn.WriteJson(&Status{
			RAM:  0.0,
			VRAM: vram,
		}); err != nil {
			fmt.Println(err)
			conn.Close()
		}
	}
}
