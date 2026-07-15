package handlers
 
import (
	db "forum/database"
	"forum/helpers"
	"net/http"
)
 
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SendJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
 
	cookie, err := r.Cookie("Form_Token")
	if err != nil {
		helpers.SendJSON(w, http.StatusUnauthorized, "No session found")
		return
	}
 
	_, err = db.DB.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
 
	http.SetCookie(w, &http.Cookie{
		Name:   "Form_Token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
 
	helpers.SendJSON(w, http.StatusOK, "Logout successful")
}
 
