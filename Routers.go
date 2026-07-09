package main

import "net/http"

var Templates = 

func Routers() *http.ServeMux {
	mux := http.NewServeMux()

	// Define your routes here
	mux.HandleFunc("/", SignIn)

	return mux

}