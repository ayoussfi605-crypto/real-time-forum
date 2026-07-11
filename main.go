package main

import (
	db "forum/backend/database"
	"forum/backend/routes"
	"log"
	"net/http"
)

func main() {

	//connect to the database first
	Conn, err := db.ConnToDb()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer Conn.Close() // Close the database connection when the main function exits

	err = db.CreateTables(Conn)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	mux := routes.Routers() // Create a new router

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
