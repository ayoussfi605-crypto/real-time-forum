package handlers

import (
	"encoding/json"
	db "forum/database"
	"forum/helpers"
	"forum/middlewares"
	"net/http"
)

type MeResponse struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var nickname string
	err := db.DB.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	// say from browser the coming response is json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(MeResponse{ID: userID, Nickname: nickname})
}
