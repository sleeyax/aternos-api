package aternos_api

type ServerInfo struct {
	// Server status (Online, Offline).
	Status string

	// IP address or hostname.
	Address string

	// Installed server software.
	Software string

	// Current supported Minecraft version.
	Version string
}
