package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "forum/database"
	"forum/helpers"
	"forum/middlewares"
)

// What we expect from the frontend when adding a comment
type CreateCommentRequest struct {
	Content string `json:"content"`
}

// POST /api/posts/{id}/comments — adds a comment to a post
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SendJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the logged-in user's ID from the auth middleware
	userID := r.Context().Value(middlewares.UserIDKey).(int)

	// Get the post ID from the URL (e.g. /api/posts/5/comments → "5")
	idStr := r.PathValue("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.SendJSON(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Read the comment content from the request body
	var input CreateCommentRequest
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		helpers.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Comment must not be empty
	if input.Content == "" {
		helpers.SendJSON(w, http.StatusBadRequest, "Comment cannot be empty")
		return
	}

	// Insert the comment into the database
	_, err = db.DB.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		postID, userID, input.Content,
	)
	fmt.Println("com id")
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not add comment")
		return
	}

	helpers.SendJSON(w, http.StatusCreated, "Comment added successfully")
}
