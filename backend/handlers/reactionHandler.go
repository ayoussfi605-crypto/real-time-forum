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

	userID := r.Context().Value(middlewares.UserIDKey).(int)

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Invalid post ID",
		)
		return
	}

	reaction, err := db.GetReaction(postID, userID)
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusInternalServerError,
			"Failed to get reaction",
		)
		return
	}
	// user doesn't have reaction
	userReaction := req.Reaction

	if reaction == nil {
		err := db.CreateReaction(postID, userID, req.Reaction)
		if err != nil {
			helpers.SendJSON(
				w,
				http.StatusInternalServerError,
				"Failed to create reaction",
			)
			return
		}
	} else {
		if reaction.Reaction == req.Reaction {
			err := db.DeleteReaction(postID, userID)
			if err != nil {
				helpers.SendJSON(
					w,
					http.StatusInternalServerError,
					"Failed to delete reaction",
				)
				return
			}

			userReaction = ""
		} else {
			err := db.UpdateReaction(postID, userID, req.Reaction)
			if err != nil {
				helpers.SendJSON(
					w,
					http.StatusInternalServerError,
					"Failed to update reaction",
				)
				return
			}
		}
	}

	stats, err := db.CountReactions(postID)
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusInternalServerError,
			"Failed to count reactions",
		)
		return
	}

	helpers.SendJSON(
		w,
		http.StatusOK,
		"Reaction updated",
		map[string]interface{}{
			"likes":        stats.Likes,
			"dislikes":     stats.Dislikes,
			"userReaction": userReaction,
		},
	)
}
