package utils

import (
	"time"
)

//Process is struct
type Process struct {
	ID           string
	UseVram      float32
	Status       string
	Filename     string
	StartDate    *time.Time
	CompleteDate *time.Time
}

// BroadcastProcess データベースの情報を格納する
var BroadcastProcess = make(chan []Process)
