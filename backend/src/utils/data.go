package utils

type CalcServers struct {
	IP     string `db:"ip"`
	Port   string `db:"port"`
	Status string `db:"status"`
}
