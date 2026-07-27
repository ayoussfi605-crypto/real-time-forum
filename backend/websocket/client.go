package ws

import (
	"encoding/json"
	"log"
)

type IncomingMessage struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

func (c *Client) ReadPump(hub *Hub) {
	defer hub.RemoveClient(c.userID)
	defer c.conn.Close()

	for {
		var msg IncomingMessage

		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println("invalid message:", err)
			continue
		}
		// Never trust the client for the sender ID
		msg.SenderID = c.userID
		log.Printf("Message received: %+v\n", msg)

		hub.HandleMessage(msg)
	}
}