package main

import (
	"log"
	"net/http"
)

func main() {

	//connect to the database first
	Conn, err := ConnToDb()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer Conn.Close() // Close the database connection when the main function exits

	mux := Routers() // Create a new router

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}