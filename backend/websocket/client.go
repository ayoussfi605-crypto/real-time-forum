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

// Client connect
//       │
//       ▼
// ServeWs()
//       │
//       ▼
// Upgrade HTTP -> WebSocket
//       │
//       ▼
// Create Client
//       │
//       ▼
// AddClient()
//       │
//       ▼
// ReadPump()
//       │
//       ▼
// ReadMessage()
//       │
//       ▼
// json.Unmarshal()
//       │
//       ▼
// msg.SenderID = c.userID
//       │
//       ▼
// HandleMessage()
//       │
//       ├───────────────┐
//       │               │
//       ▼               ▼
//  Save SQLite     Find Receiver
//                       │
//                       ▼
//              clients[ReceiverID]
//                       │
//              ┌────────┴─────────┐
//              │                  │
//           Online            Offline
//              │                  │
//              ▼                  ▼
//        WriteJSON()      Do Nothing


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