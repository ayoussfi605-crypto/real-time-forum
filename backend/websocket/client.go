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

func (c *Client) SafeWriteJSON(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conn.SetWriteDeadline(time.Now().Add(writeWait))

	return c.conn.WriteJSON(v)
}