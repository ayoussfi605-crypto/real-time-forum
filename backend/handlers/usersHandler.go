package handlers

import (
	"encoding/json"
	db "forum/database"
	"forum/middlewares"
	"log"
	"net/http"
)

type ChatUser struct {
	ID          int    `json:"id"`
	Nickname    string `json:"nickname"`
	UnreadCount int    `json:"unread_count"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	query := `
		SELECT u.id, u.nickname,
		       SUM(CASE WHEN m.receiver_id = ? AND m.is_read = 0 THEN 1 ELSE 0 END) as unread_count
		FROM users u
		LEFT JOIN messages m ON (m.sender_id = u.id AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = u.id)
		WHERE u.id != ?
		GROUP BY u.id, u.nickname
		ORDER BY MAX(m.created_at) DESC, u.nickname ASC
	`

	rows, err := db.DB.Query(query, userID, userID, userID, userID)
	if err != nil {
		log.Println("Error fetching users:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []ChatUser
	for rows.Next() {
		var u ChatUser
		if err := rows.Scan(&u.ID, &u.Nickname, &u.UnreadCount); err != nil {
			log.Println("Error scanning user:", err)
			continue
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating users:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
