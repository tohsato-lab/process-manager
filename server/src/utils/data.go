package utils

//Process is struct
type Process struct {
	ID           string
	Status       string
	Filename     string
	StartDate    string
	CompleteDate string
	TargetFile   string
	EnvName      string
	Comment      string
	InTrash      bool
}

type DirectoryInfo struct {
	Name  string
	IsDir bool
}

// BroadcastProcess データベースの情報を格納する
var BroadcastProcess = make(chan []Process)
