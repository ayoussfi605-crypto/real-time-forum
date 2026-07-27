package routes

import (
	"net/http"

	"forum/handlers"
	"forum/middlewares"
	ws "forum/websocket"
)

func InitRoutes(hub *ws.Hub) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("../frontend")))
	mux.HandleFunc("/me", middlewares.AuthMiddleware(handlers.MeHandler))
	mux.HandleFunc("/ws", middlewares.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	}))
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	return mux
}
