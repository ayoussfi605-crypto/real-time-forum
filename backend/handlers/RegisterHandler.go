package handlers

import (
	"encoding/json"
	db "forum/database"
	"forum/helper"
	"net/http"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Nickname  string `json:"nickname"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Age       uint   `json:"age"`
    Gender    string `json:"gender"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		helper.SendJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helper.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	input.Nickname = strings.TrimSpace(input.Nickname)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.TrimSpace(input.Email)

	helper.ValidateRegisterInput(input.Nickname, input.FirstName , input.LastName, input.Age, input.Gender, input.Email, w)


	if len(input.Password) < 8 || len(input.Password) > 50 {
		helper.SendJSON(w, http.StatusBadRequest, "Password must be between 8 and 50 characters long")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, "Error occurred while hashing password")
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (nickname, first_name, last_name, age, gender, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		input.Nickname, input.FirstName, input.LastName, input.Age, input.Gender, input.Email, hashedPassword)
		if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.nickname") {
				helper.SendJSON(w, http.StatusConflict, "Duplicated nickname. Please enter a valid nickname.")
			} else if strings.Contains(err.Error(), "users.email") {
				helper.SendJSON(w, http.StatusConflict, "Duplicated Email. Enter a valid email.")
			} else {
				helper.SendJSON(w, http.StatusConflict, "Duplicated data")
			}
			return 
		}
 
		helper.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
		
	w.WriteHeader(http.StatusCreated)
	helper.SendJSON(w, http.StatusCreated, "Registration successful")


}	