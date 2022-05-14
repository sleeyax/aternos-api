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
		close(w.Message)
	}()

	for {
		msgType, rawMsg, err := w.conn.ReadMessage()

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println(fmt.Sprintf("reader: server closed the connection unexpectedly [%s]", err))
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
