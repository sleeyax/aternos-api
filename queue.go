package aternos_api

// QueueReduction is a simplified version of Queue that the server uses to notify queue reduction (queue_reduced).
type QueueReduction struct {
	// Unique number to identify the queue.
	// This is usually 0.
	Number int `json:"queue"`

	// Current total amount of people in queue.
	Total int `json:"total"`

	MaxTime int `json:"maxtime"`
}

// Queue is the current status in queue.
// This struct is part of ServerInfo, unlike QueueReduction.
type Queue struct {
	// Unique number to identify the queue.
	// This is usually 0.
	Number int `json:"queue"`

	// Current position in queue.
	// The higher the number, the longer you have to wait until it's your turn.
	Position int `json:"position"`

	// Amount of people in queue.
	Count int `json:"count"`

	// Percentage in the queue is calculated based on (Position / Count) * 100.
	Percentage float32 `json:"percentage"`

	// Status message.
	// E.g "waiting".
	Status string `json:"pending"`

	// Time left in human readable format.
	// E.g "ca. 1 min"
	Time string `json:"time"`

	// Minutes left to wait in number format.
	Minutes int `json:"minutes"`

	Jointime int `json:"jointime"`
}
