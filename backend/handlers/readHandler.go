package handlers

import (
	db "forum/database"
	"forum/middlewares"
	ws "forum/websocket"
	"log"
	"net/http"
	"strconv"
)

func MarkMessagesReadHandler(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middlewares.UserIDKey).(int)
		
		otherUserIDStr := r.PathValue("id")
		otherUserID, err := strconv.Atoi(otherUserIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		_, err = db.DB.Exec("UPDATE messages SET is_read = 1 WHERE sender_id = ? AND receiver_id = ?", otherUserID, userID)
		if err != nil {
			log.Println("Error updating messages to read:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Broadcast to clear unread counts on other tabs
		hub.HandleMessage(ws.WSMessage{
			Type:       "messages_read",
			SenderID:   otherUserID,
			ReceiverID: userID,
		})

		w.WriteHeader(http.StatusOK)
	}
}
