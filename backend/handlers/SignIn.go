package handlers

import (
	"html/template"
	"net/http"
)

var Template = template.Must(template.ParseFiles("static/index.html"))

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render the sign-in page
		Template.Execute(w, nil)
		return
	}
	
	if r.Method == http.MethodPost {
		// Handle the sign-in form submission
		username := r.FormValue("username")
		password := r.FormValue("password")
	
		// Validate the username and password (you can implement your own logic here)
		if username == "admin" && password == "password" {
			// Successful sign-in
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	
		// Invalid credentials, show an error message
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	
	// If the request method is not GET or POST, return a 405 Method Not Allowed response
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}	
