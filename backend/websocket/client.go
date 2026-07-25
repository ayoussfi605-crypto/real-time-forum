package ws

import (
	"encoding/json"
	"log"
)

type IncomingMessage struct {
	SenderID   int
	ReceiverID int
	Content    string
}

func (c *Client) ReadPump(hub *Hub) {

}
