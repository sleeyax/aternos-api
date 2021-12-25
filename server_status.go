package aternos_api

type ServerStatus int

const (
	Offline ServerStatus = iota
	Preparing
	Loading
	Starting
	Online
	Stopping
)
