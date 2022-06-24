package controllers

import (
	"github.com/janberktold/sse"
	"log"
	"net/http"
	"time"

	"conda/utils"
)

type Status struct {
	USED float64
	VRAM float64
}

// Health ホストの情報を配信
func Health(w http.ResponseWriter, r *http.Request) {
	conn, _ := sse.Upgrade(w, r)
	for {
		totalVRAM := utils.GetTotalVRAM()
		// totalRAM := utils.GetTotalRAM()
		var vram = 0.0
		var usedRate = float64(utils.GetUsedRate())
		// var ram = 0.0
		if totalVRAM != 0 {
			vram = float64(utils.GetUsedVRAM() / totalVRAM)
		}
		/*
			if totalRAM != 0 {
				ram = float64(utils.GetUsedRAM() / totalRAM)
			}
		*/
		if err := conn.WriteJson(&Status{USED: usedRate, VRAM: vram}); err != nil {
			log.Println(err)
			conn.Close()
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		time.Sleep(1 * time.Second)
	}
}
