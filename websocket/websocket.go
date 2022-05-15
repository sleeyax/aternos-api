package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// New crates a new Websocket instance from an existing gorilla websocket connection.
func New(conn *websocket.Conn, reconnect RetryFunc) *Websocket {
	ws := &Websocket{
		readerDone:      make(chan interface{}),
		writerDone:      make(chan interface{}),
		Message:         make(chan Message),
		messageOut:      make(chan Message),
		readerReconnect: make(chan bool),
		writerReconnect: make(chan bool),
		conn:            conn,
		OnReconnect:     reconnect,
	}

	go ws.startReader()
	go ws.startWriter()

	return ws
}

// Close closes the underlying websocket connection and waits for all goroutines to exit.
func (w *Websocket) Close() error {
	close(w.writerDone)

	select {
	case <-w.readerDone:
		log.Println("reader & writer goroutines stopped")
		return nil
	case <-time.After(time.Duration(5) * time.Second):
		log.Println("timeout in closing reader goroutine, exiting by force....")
		return w.conn.Close()
	}
}

// Send sends a message over the websocket connection.
func (w *Websocket) Send(message Message) {
	w.messageOut <- message
}

// StartConsoleLogStream starts fetching the server start logs (console).
// This function should only be called once the server has been started over HTTP.
func (w *Websocket) StartConsoleLogStream() {
	w.Send(Message{
		Stream: "console",
		Type:   "start",
	})
}

// StopConsoleLogStream starts fetching the server stop logs (console).
// This function should only be called once the server has been stopped over HTTP.
func (w *Websocket) StopConsoleLogStream() {
	w.Send(Message{
		Stream: "console",
		Type:   "stop",
	})
}

// StartHeapInfoStream starts fetching information about the server heap.
// See https://www.javatpoint.com/java-heap for more information about heaps.
func (w *Websocket) StartHeapInfoStream() {
	w.Send(Message{
		Stream: "heap",
		Type:   "start",
	})
}

// StopHeapInfoStream stops fetching information about the server heap.
// See https://www.javatpoint.com/java-heap for more information about heaps.
func (w *Websocket) StopHeapInfoStream() {
	w.Send(Message{
		Stream: "heap",
		Type:   "stop",
	})
}

// StartTickStream starts streaming the current server tick count.
// See https://minecraft.fandom.com/wiki/Tick for more information about ticks.
func (w *Websocket) StartTickStream() {
	w.Send(Message{
		Stream: "tick",
		Type:   "start",
	})
}

// StopTickStream stops streaming the current server tick count.
// See https://minecraft.fandom.com/wiki/Tick for more information about ticks.
func (w *Websocket) StopTickStream() {
	w.Send(Message{
		Stream: "tick",
		Type:   "stop",
	})
}

// SendHeartBeat sends a single keep-alive request.
// The server doesn't respond to this request, but a heartbeat should be regularly sent to keep the connection alive.
func (w *Websocket) SendHeartBeat() {
	w.Send(Message{
		Type: "\xE2\x9D\xA4", // hearth emoji ❤
	})
}
