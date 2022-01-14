package aternos_api

type ServerInfo struct {
	// Server brand.
	// Usually this is just "aternos".
	Brand string `json:"brand"`

	// Status code.
	Status ServerStatus `json:"status"`

	// Time since the server info, such as the status, has been updated.
	LastChanged int `json:"change"`

	// Maximum amount of players.
	MaxPlayers int `json:"slots"`

	// Number of logged issues or problems.
	Problems int `json:"problems"`

	// Amount of active players.
	Players int `json:"players"`

	// List of active players.
	PlayerList []string `json:"playerlist"`

	Message struct {
		Text  string `json:"text"`
		Class string `json:"class"`
	} `json:"message"`

	// Dynamic IP address.
	DynIP string `json:"dynip"`

	// Whether this is a server for bedrock edition.
	IsBedrock bool `json:"bedrock"`

	Host string `json:"host"`

	Port int `json:"port"`

	HeadStarts int `json:"headstarts"`

	RAM int `json:"ram"`

	// Status label
	// E.g. online, offline.
	StatusLabel string `json:"lang"`

	// Status label but in human-readable format.
	// E.g. Starting..., Online, Offline.
	StatusLabelFormatted string `json:"label"`

	// CSS class that's applied to the status label.
	StatusLabelClass string `json:"class"`

	// Amount of time left to join the server.
	Countdown int `json:"countdown"`

	Queue Queue `json:"queue"`

	// Unique server ID.
	Id string `json:"id"`

	// Name of the server.
	Name string `json:"name"`

	// Server software.
	Software string `json:"software"`

	SoftwareId string `json:"softwareId"`

	// Software type.
	// E.g vanilla, papermc.
	SoftwareType string `json:"type"`

	// Current Minecraft version.
	Version string `json:"version"`

	IsDeprecated bool `json:"deprecated"`

	// IP address.
	IP string `json:"ip"`

	// Domain address.
	Address string `json:"displayAddress"`

	// Message Of The Day.
	MOTD string `json:"motd"`

	// Name of the class that's used to display an icon next to the status.
	// E.g fa-stop-circle.
	Icon string `json:"icon"`

	DNS struct {
		Type    string   `json:"type"`
		Domains []string `json:"domains"`
		Host    string   `json:"host"`
		Port    int      `json:"port"`
		IP      string   `json:"ip"`
	} `json:"dns"`
}
