package handlers

import (
	"encoding/json"
	db "forum/database"
	"forum/helpers"
	"forum/middlewares"
	"net/http"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var nickname string
	err := db.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"nickname": nickname})
}