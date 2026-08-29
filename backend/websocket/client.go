package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
)

type Client struct {
	userID int
	conn   *websocket.Conn
	mutex  sync.Mutex
}

// SafeWriteJSON writes a JSON message to the WebSocket connection in a thread-safe manner.
func (c *Client) SafeWriteJSON(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(writeWait))

	return c.conn.WriteJSON(v)
}