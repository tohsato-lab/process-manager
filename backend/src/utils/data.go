package utils

//Process is struct
type Process struct {
	ID           string
	ProcessName  string
	Status       string
	StartDate    string
	CompleteDate string
	TargetFile   string
	EnvName      string
	Comment      string
	InTrash      bool
	IsHome       bool
	ServerIP     string
}

type Servers struct {
	IP     string
	Port   string
	Status string
}

type DirectoryInfo struct {
	Name  string
	IsDir bool
}

// BroadcastProcesses データベースの情報を格納する
var BroadcastProcesses = make(chan []Process)
