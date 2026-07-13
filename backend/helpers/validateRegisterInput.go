package helpers

import (
	"net/http"
	"net/mail"
)

func ValidateRegisterInput(nickname string, firstName string, lastName string, age uint, gender string, email string, password string, w http.ResponseWriter) {
	// validate Gender
	if gender != "Male" && gender != "Female" {
		SendJSON(w, http.StatusBadRequest, "Invalid gender. Please select 'Male' or 'Female'")
		return
	}

	// Validate the request data (you can add more validation as needed)
	if nickname == "" || firstName == "" || lastName == "" ||
		age <= 0 || age > 150 || gender == "" || email == "" || password == "" {
		SendJSON(w, http.StatusBadRequest, "Please fill all the fields")
		return
	}

	if len(nickname) < 3 || len(nickname) > 20 || len(firstName) < 3 || len(firstName) > 20 || len(lastName) < 3 || len(lastName) > 20 {
		SendJSON(w, http.StatusBadRequest, "Nickname, first name, and last name must be between 3 and 20 characters long")
		return
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if len(password) < 8 || len(password) > 50 {
		SendJSON(w, http.StatusBadRequest, "Password must be between 8 and 50 characters long")
		return
	}
}
