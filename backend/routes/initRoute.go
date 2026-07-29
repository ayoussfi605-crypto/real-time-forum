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
	// /ws (identify hub conn)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// Posts
	mux.HandleFunc("GET /api/posts", handlers.GetPostsHandler)
	mux.HandleFunc("POST /api/posts", middlewares.AuthMiddleware(handlers.CreatePostHandler))
	mux.HandleFunc("GET /api/posts/{id}", handlers.GetPostByIDHandler)

	// Comments
	mux.HandleFunc("POST /api/posts/{id}/comments", middlewares.AuthMiddleware(handlers.CreateCommentHandler))

	// Categories (needed by the "new post" form to know what's selectable)
	mux.HandleFunc("GET /api/categories", handlers.GetCategoriesHandler)

	return mux
}