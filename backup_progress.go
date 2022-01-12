package aternos_api

type BackupProgress struct {
	Id string `json:"id"`

	// Percentage the backup is done.
	Progress int `json:"progress"`

	Action string `json:"action"`

	Auto bool `json:"auto"`

	// Whether the backup is complete.
	Done bool `json:"done"`
}
