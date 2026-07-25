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
	mutex   sync.Mutex
}

func NewHub() *Hub {

	return &Hub{
		clients: make(map[int]*Client),
	}
}

func (h *Hub) AddClient(userID int, client *Client) {

	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.clients[userID] = client
}

func (h *Hub) RemoveClient(userID int) {

	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.clients, userID)
}

func (h *Hub) HandleMessage(msg IncomingMessage) {
	_, err := db.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		msg.SenderID, msg.ReceiverID, msg.Content,
	)
	if err != nil {
		log.Println(err)
	}

	h.mutex.Lock()
	client, ok := h.clients[msg.ReceiverID]
	h.mutex.Unlock()

	if ok {
		// Receiver online -> sift-lih message mباشرة
		err := client.conn.WriteJSON(msg)
		if err != nil {
			log.Println("write error:", err)
		}
	}
}
