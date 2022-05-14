package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

// Websocket wraps a gorilla websocket connection and provides a bunch of useful methods for interacting with Aternos' websocket server.
type Websocket struct {
	// Whether we are connected to the websocket server.
	// isConnected bool

	// readerDone indicates whether the reader goroutine is done processing incoming messages.
	// It's considered done when the channel is closed.
	readerDone chan interface{}

	// writerDone indicates whether the messageOut goroutine is done processing outgoing messages.
	// It's considered done when the channel is closed.
	writerDone chan interface{}

	// Received messages channel.
	Message chan Message

	// Sent messages channel.
	messageOut chan Message

	// The current websocket connection.
	conn *websocket.Conn

	mu sync.Mutex
}

// Message represents a message that can be received from the websocket server.
type Message struct {
	Stream       string `json:"stream,omitempty"`
	Type         string `json:"type,omitempty"`
	Message      string `json:"Message,omitempty"`
	MessageBytes []byte `json:"-"`
	Data         Data   `json:"data,omitempty"`
	Console      string `json:"console,omitempty"`
}

// Data represents an arbitrary payload that can either contain a string or a variable object.
type Data struct {
	// Content contains either a regular string or a serialized JSON object.
	Content      string
	ContentBytes []byte
}
