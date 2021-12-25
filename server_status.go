package aternos_api

type ServerStatus int

const (
	Offline   ServerStatus = 0
	Online    ServerStatus = 1
	Preparing ServerStatus = 10
	Starting  ServerStatus = 2
	Stopping  ServerStatus = 3
	Saving    ServerStatus = 5
	Loading   ServerStatus = 6
)
