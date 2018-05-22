package wsconnector

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
)

// WsPlayerConn ...
type WsPlayerConn struct {
	conn *websocket.Conn
}

// Receive ...
func (c *WsPlayerConn) Receive() (data []byte, err error) {

	_, message, err := c.conn.ReadMessage()

	if err != nil {

		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			log.Printf("error: %v", err)
		}
		return nil, errors.New("Connection closed")
	}
	return message, err
}

// Send ...
func (c *WsPlayerConn) Send(message []byte) error {

	w, err := c.conn.NextWriter(websocket.BinaryMessage)

	if err != nil {
		return err
	}
	w.Write(message)

	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// Close ...
func (c *WsPlayerConn) Close() error {
	return c.conn.Close()
}
