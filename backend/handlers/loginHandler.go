package handlers

import (
	"encoding/json"
	"forum/helpers"
	"net/http"
	"strings"
)

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
		helpers.GetUserByIdentifier(w, identifier)
		
	}
}
