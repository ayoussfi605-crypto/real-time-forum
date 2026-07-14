package handlers

// Login
// ↓
// Check password
// ↓
// Generate UUID token
// ↓
// Insert token into sessions table
// ↓
// Send cookie
// ↓
// Browser stores cookie
// ---------------------------
// Frontend
//  |
//  POST /login
//  |
// {
//  identifier:"ayoub",
//  password:"12345678"
// }
//  |
//  ▼
// Go
//  |
//  ▼
// Validate input
//  |
//  ▼
// Find user by email OR nickname
//  |
//  ▼
// bcrypt compare
//  |
//  ▼
// Create session
//  |
//  ▼
// Set cookie
//  |
//  ▼
// Return success

import (
	"database/sql"
	"encoding/json"
	db "forum/database"
	"forum/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Handle login logic here
		var input LoginRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			helpers.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		identifier := strings.TrimSpace(input.Identifier)
		if identifier == "" || input.Password == "" {
			helpers.SendJSON(w, http.StatusBadRequest, "Please fill all the fields")
			return
		}
		user, err := db.GetUserByIdentifier(identifier)
		if err != nil {
			if err == sql.ErrNoRows {
				helpers.SendJSON(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}
			helpers.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		// var user User
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
		if err != nil {
			helpers.SendJSON(w, http.StatusBadRequest, "Invalid credentials")
			return
		}
		generatedToken := uuid.New().String()

		_, err = db.DB.Exec("INSERT INTO session (user_id, token, expiration_date) VALUES(? , ? , ?)", user.ID, generatedToken, time.Now().Add(24*time.Hour))
		if err != nil {
			helpers.SendJSON(w, http.StatusInternalServerError, "Could not create session")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "Form_Token",
			Value:    generatedToken,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteStrictMode,
			MaxAge:   24 * 60 * 60,
		})
		
		helpers.SendJSON(w, http.StatusCreated, "succes")

	}
}
