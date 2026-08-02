package ws

import (
	"encoding/json"
	"log"
)

type IncomingData struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

type IncomingMessage struct {
	EventType string       `json:"event_type"`
	Data      IncomingData `json:"data"`
}

func (c *Client) ReadPump(hub *Hub) {
	defer hub.RemoveClient(c.userID, c) 
	defer c.conn.Close()

	c.conn.SetReadLimit(maxMsgSize)

	for {
		var msg IncomingMessage

		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		// Never trust the client for the sender ID.
		// Always use the authenticated user's ID.
		msg.Data.SenderID = c.userID

		log.Printf("Received %q: %+v\n", msg.EventType, msg.Data)

		hub.HandleMessage(msg)
	}
}
