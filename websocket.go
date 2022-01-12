package aternos_api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sleeyax/aternos-api/internal/tlsadapter"
	httpx "github.com/useflyent/fhttp"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	// InterruptSignal is signalled when a user quits the program by pressing CTRL + C on the keyboard.
	InterruptSignal chan os.Signal // TODO: also check for interrupts of the HTTP API.
)

const (
	websocketUrl = "wss://aternos.org/hermes/"
)

type Websocket struct {
	// receiverDone indicates whether the receiver goroutine is done processing incoming messages.
	// It's considered done when the channel is closed.
	receiverDone chan interface{}

	// Received messages channel.
	Message chan WebsocketMessage

	// The current websocket connection.
	conn *websocket.Conn
}

func (w *Websocket) init() {
	w.receiverDone = make(chan interface{})
	w.Message = make(chan WebsocketMessage)
	go w.startReceiver()
	go w.startSender()
}

// Close closes the websocket connection.
func (w *Websocket) Close() error {
	return w.conn.Close()
}

// Send sends a message over the websocket connection.
func (w *Websocket) Send(message WebsocketMessage) error {
	return w.conn.WriteJSON(message)
}

// StartConsoleLogStream starts fetching the server start logs (console).
// This function should only be called once the server has been started over HTTP.
func (w *Websocket) StartConsoleLogStream() error {
	return w.Send(WebsocketMessage{
		Stream: "console",
		Type:   "start",
	})
}

// StopConsoleLogStream starts fetching the server stop logs (console).
// This function should only be called once the server has been stopped over HTTP.
func (w *Websocket) StopConsoleLogStream() error {
	return w.Send(WebsocketMessage{
		Stream: "console",
		Type:   "stop",
	})
}

// StartHeapInfoStream starts fetching information about the server heap.
func (w *Websocket) StartHeapInfoStream() error {
	return w.Send(WebsocketMessage{
		Stream: "heap",
		Type:   "start",
	})
}

// StopHeapInfoStream stops fetching information about the server heap.
func (w *Websocket) StopHeapInfoStream() error {
	return w.Send(WebsocketMessage{
		Stream: "heap",
		Type:   "stop",
	})
}

// StartTickStream starts streaming the time left for someone to join the server.
// Once the timer runs out and no one has joined the server, it will be stopped automatically.
func (w *Websocket) StartTickStream() error {
	return w.Send(WebsocketMessage{
		Stream: "tick",
		Type:   "start",
	})
}

// StopTickStream stops streaming the time left for someone to join the server.
// Once the timer runs out and no one has joined the server, it will be stopped automatically.
func (w *Websocket) StopTickStream() error {
	return w.Send(WebsocketMessage{
		Stream: "tick",
		Type:   "stop",
	})
}

// SendHeartBeat sends a single keep-alive request.
// The server doesn't respond to this request, but a heartbeat should be regularly sent to keep the connection alive.
func (w *Websocket) SendHeartBeat() error {
	return w.Send(WebsocketMessage{
		Type: "\xE2\x9D\xA4", // hearth emoji ❤
	})
}

// startReceiver starts listening for incoming messages.
func (w *Websocket) startReceiver() {
	defer close(w.receiverDone)

	for {
		msgType, rawMsg, err := w.conn.ReadMessage()

		if err != nil {
			log.Println("Error in receiver: ", err)
			return
		}

		switch msgType {
		case websocket.TextMessage:
			var msg WebsocketMessage
			if err = json.Unmarshal(rawMsg, &msg); err != nil {
				log.Println("Receiver failed to parse msg: ", err)
				return
			}

			if msg.Message != "" {
				msg.MessageBytes = []byte(msg.Message)
			}

			w.Message <- msg
		case websocket.CloseMessage:
			return
		default:
			log.Printf("Unknown message received: %s\n", rawMsg)
		}
	}
}

func (w *Websocket) startSender() {
	for {
		select {
		case <-InterruptSignal:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT InterruptSignal signal. Closing all pending connections")

			// Try to close the websocket connection.
			// If for some reason it fails, just quit already because the receiver won't know the connection is ending.
			if err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("Sender failed to tell the websocket server to quit: ", err)
				return
			}

			select {
			case <-w.receiverDone:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(3) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}

			return
		}
	}
}

// ConnectWebSocket connects to the Aternos websockets server.
func (api *Api) ConnectWebSocket() (*Websocket, error) {
	headers := api.client.Options.Headers.Clone()
	headers.Set("accept", "*/*")
	headers.Set("cache-control", "no-cache")
	headers.Set("host", "aternos.org")
	headers.Set("origin", api.client.Options.PrefixURL)
	headers.Del(httpx.HeaderOrderKey)

	dialer := websocket.Dialer{
		Proxy: func(request *http.Request) (*url.URL, error) {
			return api.options.Proxy, nil
		},
		HandshakeTimeout:  30 * time.Second,
		EnableCompression: true,
		NetDialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			adapter := (api.client.Options.Adapter).(*tlsadapter.TLSAdapter)
			return adapter.DialTLSContext(ctx, network, addr)
		},
		Jar: api.client.Options.CookieJar,
	}

	conn, _, err := dialer.Dial(websocketUrl, headers)
	if err != nil {
		return nil, err
	}

	wss := &Websocket{conn: conn}
	wss.init()

	return wss, nil
}
