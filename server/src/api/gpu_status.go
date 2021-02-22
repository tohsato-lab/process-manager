package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"../utils"
)

// GPUSstatus GPUの情報を配信
func GPUSstatus(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 1秒おきにデータを流す
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				if utils.GetTotalVRAM() != 0 {
					fmt.Fprintf(w, "data: %f\n\n", utils.GetUsedVRAM()/utils.GetTotalVRAM())
				} else {
					fmt.Fprintf(w, "data: %f\n\n", 0.0)
				}
				flusher.Flush()
			}
		}
	}()
	<-r.Context().Done()
	log.Println("コネクションが閉じました")
}
