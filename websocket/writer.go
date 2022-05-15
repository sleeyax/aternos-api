package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func (w *Websocket) startWriter() {
	defer close(w.messageOut)
	defer close(w.writerReconnect)

	for {
		select {
		case message := <-w.messageOut:
			if err := w.conn.WriteJSON(message); err != nil {
				log.Println("writer: failed to write message: ", err)
			}
		case <-w.writerDone: // TODO: fix this
			// Try to tell the server that we want to close the connection.
			if err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("writer: failed to send close message: ", err)
			}

			return
		case <-w.writerReconnect:
			w.conn.Close()

			var err error

			if w.conn, err = w.OnReconnect(); err != nil {
				log.Println(fmt.Sprintf("reader: failed to reconnect [%s]", err))
				return
			}

			w.readerReconnect <- true
		}
	}
}
