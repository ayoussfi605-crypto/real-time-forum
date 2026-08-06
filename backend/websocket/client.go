package ws

import (
	"encoding/json"
	"log"
	"time"
)

type WSMessage struct {
	Type        string `json:"type"`
	SenderID    int    `json:"sender_id,omitempty"`
	ReceiverID  int    `json:"receiver_id,omitempty"`
	Content     string `json:"content,omitempty"`
	UserID      int    `json:"user_id,omitempty"` // For status
	Online      bool   `json:"online,omitempty"`  // For status
	CreatedAt   string `json:"created_at,omitempty"`
	OnlineUsers []int  `json:"online_users,omitempty"` // For initial status
}

func (c *Client) ReadPump(hub *Hub) {
	defer hub.RemoveClient(c.userID)
	defer c.conn.Close()

	for {
		var msg WSMessage

		_, p, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println("invalid message:", err)
			continue
		}
		
		if msg.Type == "chat_message" {
			msg.SenderID = c.userID
			msg.CreatedAt = time.Now().Format(time.RFC3339)
			hub.HandleMessage(msg)
		}
	}
}