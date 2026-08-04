package ws

import (
	db "forum/database"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	maxMsgSize = 4096
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

type Client struct {
	userID int
	conn   *websocket.Conn
	mutex  sync.Mutex 
}


func (c *Client) SafeWriteJSON(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// deadline if the write > 10 sc stop exp if the user close chat
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteJSON(v)
}


type Hub struct {
	clients map[int]map[*Client]bool 
	mutex   sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int]map[*Client]bool),
	}
}

func (h *Hub) AddClient(userID int, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// check if user_id He has map if not make it
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]bool)
	}
	h.clients[userID][client] = true
}

func (h *Hub) RemoveClient(userID int,client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	// Remove connection
	delete(h.clients[userID], client)
	// remove map if the not have any tab or conn 
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID) 
	}
}

func (h *Hub) HandleMessage(msg IncomingMessage) {
	_, err := db.DB.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
		msg.Data.SenderID, msg.Data.ReceiverID, msg.Data.Content,
	)
	if err != nil {
		log.Println(err)
		return 
	}

	h.mutex.RLock()
	// store ReceiverID in receiverClients if user is oneline (user existing in map h.clients)
	receiverClients := h.clients[msg.Data.ReceiverID]
	h.mutex.RUnlock()

	// Brodcast message if user oneline
	for c := range receiverClients {	
		if err := c.SafeWriteJSON(msg); err != nil { 
			log.Println("write error:", err)
		}
	}
}
