package routes

import (
	"forum/backend/handlers"
	"net/http"
)

func Routers() *http.ServeMux {
	mux := http.NewServeMux()

	// Define your routes here
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", handlers.SignIn)

	return mux

}
