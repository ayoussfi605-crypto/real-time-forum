package helpers

import (
	"database/sql"
	db "forum/database"
	"net/http"
)

type User struct {
	ID       int
	Nickname string
	Password string
}

func GetUserByIdentifier(w http.ResponseWriter, identifier string) {
	var err error
	var user User
	err = db.DB.QueryRow("SELECT id, Nickname, Password FROM users WHERE Nickname = ? or Email = ?", identifier, identifier).Scan(&user.ID, &user.Nickname, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			SendJSON(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	SendJSON(w, http.StatusAccepted, "login succecful") 
}
