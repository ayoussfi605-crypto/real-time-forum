package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	db "forum/database"
	"forum/helpers"
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
		helpers.SendJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var input RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helpers.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	input.Nickname = strings.TrimSpace(input.Nickname)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.TrimSpace(input.Email)
	
	Isvalid:= helpers.ValidateRegisterInput(input.Nickname, input.FirstName, input.LastName, input.Age, input.Gender, input.Email, input.Password, w)
	if !Isvalid{
		return 
	}
	hashedPassword, err := helpers.HashPassword(input.Password, w)
	if err != nil {
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (nickname, first_name, last_name, age, gender, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		input.Nickname, input.FirstName, input.LastName, input.Age, input.Gender, input.Email, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "users.nickname") {
				helpers.SendJSON(w, http.StatusConflict, "Duplicated nickname. Please enter a valid nickname.")
			} else if strings.Contains(err.Error(), "users.email") {
				helpers.SendJSON(w, http.StatusConflict, "Duplicated Email. Enter a valid email.")
			} else {
				helpers.SendJSON(w, http.StatusConflict, "Duplicated data")
			}
			return
		}

		helpers.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	helpers.SendJSON(w, http.StatusCreated, "Registration successful")
}
