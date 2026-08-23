package ws

import "testing"

func TestClientsForUsersIncludesSenderAndReceiverTabs(t *testing.T) {
	client1Tab1 := &Client{userID: 1}
	client1Tab2 := &Client{userID: 1}
	client2Tab1 := &Client{userID: 2}

	h := &Hub{
		clients: map[int]map[*Client]bool{
			1: {
				client1Tab1: true,
				client1Tab2: true,
			},
			2: {
				client2Tab1: true,
			},
		},
	}

	clients := h.clientsForUsers(1, 2)

	if len(clients) != 3 {
		t.Fatalf("expected 3 connected clients, got %d", len(clients))
	}

	seen := map[*Client]bool{}
	for _, client := range clients {
		seen[client] = true
	}

	if !seen[client1Tab1] || !seen[client1Tab2] || !seen[client2Tab1] {
		t.Fatal("expected sender and receiver tabs to be included in outgoing delivery")
	}
}
