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
	
	// Posts
	mux.HandleFunc("GET /api/posts", handlers.GetPostsHandler)
	mux.HandleFunc("POST /api/posts", middlewares.AuthMiddleware(handlers.CreatePostHandler))
	
	mux.HandleFunc("GET /api/posts/{id}", middlewares.AuthMiddleware(handlers.GetPostByIDHandler),)
	
	// Like / Dislike
	mux.HandleFunc("POST /api/posts/{id}/reaction",middlewares.AuthMiddleware(handlers.HandleReaction),)
	
	// Comments
	mux.HandleFunc("POST /api/posts/{id}/comments", middlewares.AuthMiddleware(handlers.CreateCommentHandler))

	// Categories (needed by the "new post" form)
	mux.HandleFunc("GET /api/categories", handlers.GetCategoriesHandler)

	// Chat
	mux.HandleFunc("GET /api/users", middlewares.AuthMiddleware(handlers.GetUsersHandler))
	mux.HandleFunc("GET /api/messages/{id}", middlewares.AuthMiddleware(handlers.GetMessagesHandler))
	mux.HandleFunc("POST /api/messages/{id}/read", middlewares.AuthMiddleware(handlers.MarkMessagesReadHandler(hub)))

	return mux
}
