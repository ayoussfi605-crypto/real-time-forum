package ws

import (
	"forum/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	//Check if the request Origin header is acceptable and if the origin host is not equal to request Host header return false.
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:8080"
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.UserIDKey).(int)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		conn:   conn,
		userID: userID,
	}

	hub.AddClient(userID, client)

	client.ReadPump(hub)
}
