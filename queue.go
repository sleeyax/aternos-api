package aternos_api

type Queue struct {
	Number  int `json:"queue"`
	Total   int `json:"total"`
	MaxTime int `json:"maxtime"`
}
