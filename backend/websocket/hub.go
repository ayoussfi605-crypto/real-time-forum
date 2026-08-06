package ws

import (
	db "forum/database"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Ayoub
//    │
//    │ ws.send()
//    ▼
// ReadPump()
//    │
//    ▼
// IncomingMessage
//    │
//    ▼
// HandleMessage()
//    │
//    ├──────────────┐
//    │              │
//    ▼              ▼
// Save DB      Find Receiver
//                   │
//                   │
//            clients[ReceiverID]
//                   │
//           ┌───────┴────────┐
//           │                │
//       Online           Offline
//           │                │
//           ▼                ▼
// WriteJSON()         Do nothing

// Client kaymtl user wa7d connecté
type Client struct {
	userID int
	conn   *websocket.Conn
}

// Hub howa dfter l3anawin: kayhtafed b ga3 les clients connectés
type Hub struct {
	clients map[int]*Client //(userID -> Client)
	mutex   sync.RWMutex
}

func NewHub() *Hub {

	return &Hub{
		clients: make(map[int]*Client),
	}
}

func (h *Hub) BroadcastStatus(userID int, online bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	msg := WSMessage{
		Type:   "user_status",
		UserID: userID,
		Online: online,
	}

	for id, client := range h.clients {
		if id != userID {
			err := client.conn.WriteJSON(msg)
			if err != nil {
				log.Println("broadcast error:", err)
			}
		}
	}
}

func (h *Hub) AddClient(userID int, client *Client) {
	h.mutex.Lock()
	h.clients[userID] = client
	h.mutex.Unlock()

	h.BroadcastStatus(userID, true)
}

func (h *Hub) RemoveClient(userID int) {
	h.mutex.Lock()
	delete(h.clients, userID)
	h.mutex.Unlock()

	h.BroadcastStatus(userID, false)
}

func (h *Hub) HandleMessage(msg WSMessage) {
	_, err := db.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		msg.SenderID, msg.ReceiverID, msg.Content,
	)
	if err != nil {
		log.Println("db insert error:", err)
		return
	}

	h.mutex.RLock()
	client, ok := h.clients[msg.ReceiverID]
	h.mutex.RUnlock()

	if ok {
		// Receiver online -> sift-lih message mباشرة
		err := client.conn.WriteJSON(msg)
		if err != nil {
			log.Println("write error:", err)
		}
	}
}
