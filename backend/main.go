package main

import (
	db "forum/database"
	"forum/routes"
	ws "forum/websocket"
	"log"
	"net/http"
)

func main() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	hub := ws.NewHub()
	mux := routes.InitRoutes(hub)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
