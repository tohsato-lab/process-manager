package utils

//Process is struct
type Process struct {
	ID           string
	UseVram      float32
	Status       string
	Filename     string
	StartDate    string
	CompleteDate string
	TargetFile   string
	EnvName      string
	ExecCount    int32
}

// BroadcastProcess データベースの情報を格納する
var BroadcastProcess = make(chan []Process)
