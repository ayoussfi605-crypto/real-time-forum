package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "forum/database"
	"forum/helpers"
	"forum/middlewares"
)

type ReactionRequest struct {
	Reaction string `json:"reaction"`
}

func HandleReaction(w http.ResponseWriter, r *http.Request) {
	// 1. Read JSON body
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

	// 2. Validate reaction
	if req.Reaction != "like" && req.Reaction != "dislike" {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Reaction must be like or dislike",
		)
		return
	}

	// 3. Get authenticated user ID
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	// 4. Get post ID from URL
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Invalid post ID",
		)
		return
	}

	db.GetReaction(postID, userID)
}
