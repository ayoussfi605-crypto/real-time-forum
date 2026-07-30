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
        //      Sara Browser

// upgrader: l config li khass l "tarjama" mn HTTP l WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 1. Jib userID (Hint: rah mkhzn f context, bhal f MeHandler:
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	// 2. Dir l "upgrade" mn HTTP l WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 3. Khl9 Client jdid b (userID, conn)
	client := &Client{
		conn:   conn,
		userID: userID,
	}
	// 4. Zido l Hub (AddClient)
	hub.AddClient(userID, client)
	// 5. TODO (khotwa jaya): bda reading loop bach ne9raw messages jayin mn had client
	client.ReadPump(hub)
}
