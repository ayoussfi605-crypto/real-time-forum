package helper

import (
	"net/http"
	"net/mail"
)


func ValidateRegisterInput(nickname string, firstName string, lastName string, age uint, gender string, email string, w http.ResponseWriter) {	


	//validate Gender
	if gender != "Male" && gender != "Female" {
		helper.SendJSON(w, http.StatusBadRequest, "Invalid gender. Please select 'Male' or 'Female'")
		return
	}

	// Validate the request data (you can add more validation as needed)
	if nickname == "" || firstName == "" || lastName == "" ||
		age <= 0 || gender == "" || email == "" {
		helper.SendJSON(w, http.StatusBadRequest, "Please fill all the fields")
		return
	}

	if len(nickname) < 3 || len(nickname) > 20 {
		helper.SendJSON(w, http.StatusBadRequest, "Nickname must be between 3 and 20 characters long")
		return
	}
	_, err := mail.ParseAddress(email)
		if err != nil {
			helper.SendJSON(w, http.StatusBadRequest, "Invalid email address")
			return
	}
}