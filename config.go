package aternos_api

type Config struct {
	// Aternos token.
	// This can be found at the end of the URL when starting the server.
	Token string

	// Authorization cookies.
	//
	// Each cookie must be formatted as NAME=VALUE.
	//
	// It must include at least ATERNOS_SESSION, ATERNOS_SEC_xxx & ATERNOS_SERVER.
	//
	// It's recommended to also specify ATERNOS_LANGUAGE.
	Cookies []string
}
