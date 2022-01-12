package aternos_api

type WebsocketMessage struct {
	Stream       string `json:"stream,omitempty"`
	Type         string `json:"type,omitempty"`
	Message      string `json:"Message,omitempty"`
	MessageBytes []byte `json:"-"`
	Data         Data   `json:"data,omitempty"`
	Console      string `json:"console,omitempty"`
}
