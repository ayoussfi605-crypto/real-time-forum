package ws

import (
	"encoding/json"
	"forum/middlewares"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type        string `json:"type"`
	SenderID    int    `json:"sender_id,omitempty"`
	ReceiverID  int    `json:"receiver_id,omitempty"`
	Content     string `json:"content,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
	Online      bool   `json:"online,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	OnlineUsers []int  `json:"online_users,omitempty"`
}
// check if the origin is allowed to connect
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &Client{
		conn:   conn,
		userID: userID,
	}

	hub.AddClient(userID, client)

	hub.mutex.RLock()

	onlineUsers := make([]int, 0, len(hub.clients))

	// Collect all unique user IDs of online users
	for id := range hub.clients {
		onlineUsers = append(onlineUsers, id)
	}

	hub.mutex.RUnlock()

	// Send the initial status message to the newly connected client
	err = client.SafeWriteJSON(WSMessage{
		Type:        "initial_status",
		OnlineUsers: onlineUsers,
	})

	if err != nil {
		log.Println("initial status error:", err)
	}

	client.ReadPump(hub)
}

// ReadPump reads messages from the WebSocket connection and processes them.
func (c *Client) ReadPump(hub *Hub) {
	// Ensure the client is removed from the hub and the connection is closed when this function exits
	defer hub.RemoveClient(c.userID, c)
	defer c.conn.Close()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var msg WSMessage

		err = json.Unmarshal(payload, &msg)

		if err != nil {
			log.Println("message unmarshal error:", err)
			continue
		}

		// Validate message type and content
		if msg.Type != "chat_message" && msg.Type != "typing" && msg.Type != "stop_typing" {
			continue
		}

		// Ignore empty chat messages
		if msg.Type == "chat_message" && msg.Content == "" {
			continue
		}

		// Never trust sender_id from frontend
		msg.SenderID = c.userID

		// Cannot send message to yourself
		if msg.SenderID == msg.ReceiverID {
			continue
		}

		// Server creates timestamp
		// RFC3339 is a standard format for timestamps in JSON and is widely used in APIs.
		msg.CreatedAt = time.Now().Format(time.RFC3339)

		hub.HandleMessage(msg)
	}
}