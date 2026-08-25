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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	conn.SetReadLimit(maxMsgSize)

	client := &Client{
		conn:   conn,
		userID: userID,
	}

	hub.AddClient(userID, client)

	hub.mutex.RLock()

	onlineUsers := make([]int, 0, len(hub.clients))

	for id := range hub.clients {
		onlineUsers = append(onlineUsers, id)
	}

	hub.mutex.RUnlock()

	err = client.SafeWriteJSON(WSMessage{
		Type:        "initial_status",
		OnlineUsers: onlineUsers,
	})

	if err != nil {
		log.Println("initial status error:", err)
	}

	client.ReadPump(hub)
}

func (c *Client) ReadPump(hub *Hub) {

	defer hub.RemoveClient(c.userID, c)
	defer c.conn.Close()

	for {

		_, payload, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg WSMessage

		err = json.Unmarshal(payload, &msg)

		if err != nil {
			log.Println("invalid message:", err)
			continue
		}

		if msg.Type != "chat_message" && msg.Type != "typing" && msg.Type != "stop_typing" {
			continue
		}

		// Empty message only matters for chat messages
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
		msg.CreatedAt = time.Now().Format(time.RFC3339)

		hub.HandleMessage(msg)
	}
}