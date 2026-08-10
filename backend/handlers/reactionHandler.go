package handlers

import (
	"encoding/json"
	"net/http"

	"forum/helpers"
)

type ReactionRequest struct {
	Reaction string `json:"reaction"`
}

func HandleReaction(w http.ResponseWriter, r *http.Request) {
	var req ReactionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	if req.Reaction != "like" && req.Reaction != "dislike" {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Reaction must be like or dislike",
		)
		return
	}
}