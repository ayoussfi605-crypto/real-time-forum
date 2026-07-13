package handlers

import (
	"backend/helpers"
	"encoding/json"
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle login logic here
		var input RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			helpers.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		input.Email = strings.TrimSpace(input.Email)
		input.Password = strings.TrimSpace(input.Password)

		if input.Email == "" || input.Password == "" {
			helpers.SendJSON(w, http.StatusBadRequest, "Please fill all the fields")
			return
		}

		user, err := db.GetUserByEmail(input.Email)
		if err != nil {
			if err == db.ErrUserNotFound {
				helpers.SendJSON(w, http.StatusUnauthorized, "Invalid email or password")
			} else {
				helpers.SendJSON(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		if !helpers.CheckPasswordHash(input.Password, user.Password) {
			helpers.SendJSON(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		token, err := helpers.GenerateJWT(user.ID)
		if err != nil {
			helpers.SendJSON(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		helpers.SendJSON(w, http.StatusOK, map[string]string{"token": token})
	}
}
