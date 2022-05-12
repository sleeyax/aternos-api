package aternos_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sleeyax/aternos-api/internal/tlsadapter"
	httpx "github.com/useflyent/fhttp"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	websocketUrl = "wss://aternos.org/hermes/"
)

type Websocket struct {
	// Whether we are connected.
	isConnected bool

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
}

// IsConnected returns whether we are connected.
func (w *Websocket) IsConnected() bool {
	return w.isConnected
}

// Close closes the websocket connection.
func (w *Websocket) Close() error {
	// Try to tell the server that we want to close the connection.
	if err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		return fmt.Errorf("failed send close message: %e", err)
	}

	w.isConnected = false

	select {
	case <-w.receiverDone:
		log.Println("Receiver Channel Closed! Exiting....")
		return nil
	case <-time.After(time.Duration(3) * time.Second):
		log.Println("Timeout in closing receiving channel. Exiting by force....")
		return w.conn.Close()
	}
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
// See https://www.javatpoint.com/java-heap for more information about heaps.
func (w *Websocket) StartHeapInfoStream() error {
	return w.Send(WebsocketMessage{
		Stream: "heap",
		Type:   "start",
	})
}

// StopHeapInfoStream stops fetching information about the server heap.
// See https://www.javatpoint.com/java-heap for more information about heaps.
func (w *Websocket) StopHeapInfoStream() error {
	return w.Send(WebsocketMessage{
		Stream: "heap",
		Type:   "stop",
	})
}

// StartTickStream starts streaming the current server tick count.
// See https://minecraft.fandom.com/wiki/Tick for more information about ticks.
func (w *Websocket) StartTickStream() error {
	return w.Send(WebsocketMessage{
		Stream: "tick",
		Type:   "start",
	})
}

// StopTickStream stops streaming the current server tick count.
// See https://minecraft.fandom.com/wiki/Tick for more information about ticks.
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

// SendHearthBeats keeps sending keep-alive requests at a specified interval.
// If no interval is specified, a default is used.
// It's recommended to use the default value unless you have a good reason not to do so.
//
// See Websocket.SendHeartBeat for more information.
func (w *Websocket) SendHearthBeats(ctx context.Context, duration ...time.Duration) {
	d := time.Millisecond * 49000
	if len(duration) > 0 {
		d = duration[0]
	}

	ticker := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			w.SendHeartBeat()
		}
	}
}

// startReceiver starts listening for incoming messages.
func (w *Websocket) startReceiver() {
	defer close(w.receiverDone)
	defer close(w.Message)

	for {
		msgType, rawMsg, err := w.conn.ReadMessage()

		if err != nil {
			closeErr, ok := err.(*websocket.CloseError)

			if !ok || closeErr.Code != websocket.CloseNormalClosure {
				log.Println("Unknown error in receiver: ", err)
			}

			w.isConnected = !ok

			return
		}

		switch msgType {
		case websocket.TextMessage:
			var msg WebsocketMessage
			if err = json.Unmarshal(rawMsg, &msg); err != nil {
				log.Println("Receiver failed to parse msg: ", err)
				break
			}

			if msg.Message != "" {
				msg.MessageBytes = []byte(msg.Message)
			}

			w.Message <- msg
		case websocket.CloseMessage:
			w.isConnected = false
			return
		default:
			log.Printf("Unknown message received: %s\n", rawMsg)
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
		Proxy:             http.ProxyURL(api.options.Proxy),
		HandshakeTimeout:  30 * time.Second,
		EnableCompression: true,
		ProxyTLSConnection: func(ctx context.Context, proxyConn net.Conn) (net.Conn, error) {
			adapter := (api.client.Options.Adapter).(*tlsadapter.TLSAdapter)
			return adapter.ConnectTLSContext(ctx, proxyConn)
		},
		Jar: api.client.Options.CookieJar,
	}

	conn, _, err := dialer.Dial(websocketUrl, headers)
	if err != nil {
		return nil, err
	}

	wss := &Websocket{conn: conn, isConnected: true}
	wss.init()

	return wss, nil
}
