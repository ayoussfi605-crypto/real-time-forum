package routes

import (
	"forum/handlers"
	"net/http"
)

func InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", handlers.RegisterHandler)
	// mux.HandleFunc("/api/login", handlers.LoginHandler)

	return mux
}