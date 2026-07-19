package helpers

import (
	"net/http"
	"net/mail"
)

func ValidateRegisterInput(nickname string, firstName string, lastName string, age uint, gender string, email string, password string, w http.ResponseWriter)  bool{
	
	if gender != "Male" && gender != "Female" {
		SendJSON(w, http.StatusBadRequest, "Invalid gender. Please select 'Male' or 'Female'")
		return false
	}

	if nickname == "" || firstName == "" || lastName == "" ||
		age <= 0 || age > 150 || email == "" || password == "" {
		SendJSON(w, http.StatusBadRequest, "Please fill all the fields")
		return false
	}

	if len(nickname) < 3 || len(nickname) > 50 || len(firstName) < 3 || len(firstName) > 50 || len(lastName) < 3 || len(lastName) > 50 {
		SendJSON(w, http.StatusBadRequest, "Nickname, first name, and last name must be between 3 and 50 characters long")
		return false
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, "Invalid email address")
		return false
	}

	if len(password) < 8 || len(password) > 50 {
		SendJSON(w, http.StatusBadRequest, "Password must be between 8 and 50 characters long")
		return false
	}
	return true
}
