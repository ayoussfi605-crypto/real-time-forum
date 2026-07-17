package handlers

import (
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

	// hna khasek function jdida li kat-encode DATA (machi ghi message)
	// ghadi n3tik l hint mn ba3d ila bghiti
}
