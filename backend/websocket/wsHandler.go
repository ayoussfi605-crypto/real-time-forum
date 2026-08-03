package ws

import (
	"forum/middlewares"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


        //      Ayoub Browser
        //           │
        //           │
        // socket.send(msg)
        //           │
        //           ▼
        // ┌──────────────────┐
        // │   Go Backend      │
        // │──────────────────│
        // │ ServeWs          │
        // │ ReadPump         │
        // │ HandleMessage    │
        // │ Save SQLite      │
        // │ Find Receiver    │
        // │ WriteJSON        │
        // └──────────────────┘
        //           │
        //           ▼
        // socket.onmessage(...)
        //           │
        //           ▼
        //      khadija el Browser

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
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
