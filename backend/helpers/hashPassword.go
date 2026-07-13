package helpers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, w http.ResponseWriter) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, "Error occurred while hashing password")
		return nil, err
	}
	return hashedPassword, nil
}
