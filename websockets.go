package aternos_api

import (
	"context"
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

const (
	websocketUrl = "wss://aternos.org/hermes/"
)

var (
	done      chan interface{}
	interrupt chan os.Signal
)

func receiveHandler(connection *websocket.Conn) {
	defer close(done)

	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		log.Printf("Received: %s\n", msg)
	}
}

// ConnectWebSocket connects to the Aternos websockets server.
func (api *AternosApi) ConnectWebSocket() error {
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
		return err
	}

	defer conn.Close()

	go receiveHandler(conn)

	// Our main loop for the client
	// We send our relevant packets here
	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}

			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}

			return nil
		}
	}
}
