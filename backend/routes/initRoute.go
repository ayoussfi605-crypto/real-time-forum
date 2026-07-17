package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middlewares"
)

func InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("../frontend")))
	mux.HandleFunc("/me", middlewares.AuthMiddleware(handlers.MeHandler))
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	// mux.HandleFunc("/api/posts", middlewares.AuthMiddleware(handlers.CreatePostHandler))
	return mux
}
