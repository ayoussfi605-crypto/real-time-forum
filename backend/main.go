package main

import (
	db "forum/database"
	"log"
	"net/http"
	"forum/routes"
)


func main() {
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := routes.InitRoutes()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

