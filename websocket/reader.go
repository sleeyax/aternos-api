package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// startReader starts listening for incoming messages.
// This method should be called in a separate goroutine and act as the only reader on the connection.
func (w *Websocket) startReader() {
	defer func() {
		close(w.readerDone)
		close(w.readerReconnect)
		close(w.Message)
	}()

	for {
		msgType, rawMsg, err := w.conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println(fmt.Sprintf("reader: server closed the connection unexpectedly [%s], reconnecting...", err))

				// tell writer to retry the connection
				w.writerReconnect <- true

				// wait for signal to read first message after reconnect
				// TODO: what if writer fails to reconnect and never notifies back? perhaps we should use a timout...
				<-w.readerReconnect

				// skip next "ready" message, so we do'nt start the server and heartbeat again
				if _, _, err = w.conn.ReadMessage(); err != nil {
					log.Println(fmt.Sprintf("reader: failed to skip next message [%s]", err))
					return
				}

				continue
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println(fmt.Sprintf("reader: server closed the connection gracefully [%s]", err))
			} else {
				log.Println(fmt.Sprintf("reader: unknown error [%s]", err))
			}

			return
		}

		switch msgType {
		case websocket.TextMessage:
			var msg Message

			if err = json.Unmarshal(rawMsg, &msg); err != nil {
				log.Println("reader: failed to parse msg: ", err)
				break
			}

			if msg.Message != "" {
				msg.MessageBytes = []byte(msg.Message)
			}

			w.Message <- msg
		default:
			log.Printf("reader: unknown message received: %s\n", rawMsg)
		}
	}
}
