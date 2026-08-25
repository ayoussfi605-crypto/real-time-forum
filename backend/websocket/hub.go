package ws

import (
	db "forum/database"
	"log"
	"sync"
)

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

	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]bool)
	}

	h.clients[userID][client] = true

	h.mutex.Unlock()

	h.BroadcastStatus(userID, true)
}

func (h *Hub) RemoveClient(userID int, client *Client) {
	h.mutex.Lock()

	if connections, ok := h.clients[userID]; ok {
		delete(connections, client)

		if len(connections) == 0 {
			delete(h.clients, userID)

			h.mutex.Unlock()

			h.BroadcastStatus(userID, false)
			return
		}
	}

	h.mutex.Unlock()
}

func (h *Hub) BroadcastStatus(userID int, online bool) {
	msg := WSMessage{
		Type:   "user_status",
		UserID: userID,
		Online: online,
	}

	h.mutex.RLock()

	var clients []*Client

	for id, connections := range h.clients {
		if id == userID {
			continue
		}

		for client := range connections {
			clients = append(clients, client)
		}
	}

	h.mutex.RUnlock()

	for _, client := range clients {
		if err := client.SafeWriteJSON(msg); err != nil {
			log.Println("broadcast error:", err)
		}
	}
}

func (h *Hub) clientsForUsers(userIDs ...int) []*Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	seen := make(map[*Client]bool)
	clients := make([]*Client, 0)

	for _, userID := range userIDs {
		for client := range h.clients[userID] {
			if !seen[client] {
				seen[client] = true
				clients = append(clients, client)
			}
		}
	}

	return clients
}

func (h *Hub) HandleMessage(msg WSMessage) {
	if msg.Type == "chat_message" {
		_, err := db.DB.Exec(
			"INSERT INTO messages (sender_id, receiver_id, content) VALUES (?, ?, ?)",
			msg.SenderID,
			msg.ReceiverID,
			msg.Content,
		)

		if err != nil {
			log.Println("db insert error:", err)
			return
		}
	}

	for _, client := range h.clientsForUsers(msg.SenderID, msg.ReceiverID) {
		if err := client.SafeWriteJSON(msg); err != nil {
			log.Println("write error:", err)
		}
	}
}
