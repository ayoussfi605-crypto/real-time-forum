package middlewares

import (
	db "forum/database"
	"forum/helpers"
	"net/http"
	"time"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("Form_Token")
		if err != nil {
			helpers.SendJSON(w, http.StatusUnauthorized, "Not authenticated")
			return
		}

		var userID int
		var expiresAt time.Time 
		err = db.DB.QueryRow(
			`SELECT user_id, expiration_date FROM sessions WHERE token = ?`,
			cookie.Value,
		).Scan(&userID, &expiresAt)
		if err != nil {
			helpers.SendJSON(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		if time.Now().After(expiresAt) {
			helpers.SendJSON(w, http.StatusUnauthorized, "Session expired")
			return
		}

		next(w, r)
	}
}