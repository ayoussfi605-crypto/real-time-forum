package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

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
