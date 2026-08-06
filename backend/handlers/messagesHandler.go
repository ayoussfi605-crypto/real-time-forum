package handlers

import (
	"encoding/json"
	"forum/middlewares"
	db "forum/database"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	otherUserIDStr := r.PathValue("id")
	otherUserID, err := strconv.Atoi(otherUserIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	query := `
		SELECT id, sender_id, receiver_id, content, created_at
		FROM messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at DESC, id DESC
		LIMIT 10 OFFSET ?
	`

	rows, err := db.DB.Query(query, userID, otherUserID, otherUserID, userID, offset)
	if err != nil {
		log.Println("Error fetching messages:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			log.Println("Error scanning message:", err)
			continue
		}
		messages = append(messages, m)
	}

	// The frontend will likely want them in chronological order (oldest first), 
	// but we fetched DESC for limit/offset. We can reverse them here or in JS.
	// Reversing here is easier for the frontend.
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
