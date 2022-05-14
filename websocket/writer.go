package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

func (w *Websocket) startWriter() {
	defer close(w.messageOut)

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
		}
	}
}
