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

type Comment struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type PostDetail struct {
	ID         int          `json:"id"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	Author     string       `json:"author"`
	CreatedAt  string       `json:"created_at"`
	Categories []string     `json:"categories"`
	Comments   []Comment    `json:"comments"`
	Reactions  ReactionInfo `json:"reactions"`
}

type ReactionInfo struct {
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	UserReaction string `json:"user_reaction"`
}

func GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Get post ID

	idStr := r.PathValue("id")

	postID, err := strconv.Atoi(idStr)
	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusBadRequest,
			"Invalid post ID",
		)
		return
	}

	// 2. Get authenticated user ID

	userID := r.Context().Value(middlewares.UserIDKey).(int)

	var post PostDetail

	// 3. Get post

	err = db.DB.QueryRow(`
		SELECT
			p.id,
			p.title,
			p.content,
			p.created_at,
			u.nickname
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = ?
	`, postID).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Author,
	)

	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusNotFound,
			"Post not found",
		)
		return
	}

	// 4. Get categories

	catRows, err := db.DB.Query(`
		SELECT c.name
		FROM categories c
		JOIN post_categories pc
			ON pc.category_id = c.id
		WHERE pc.post_id = ?
	`, postID)

	if err == nil {

		for catRows.Next() {

			var name string

			if err := catRows.Scan(&name); err != nil {
				catRows.Close()

				helpers.SendJSON(
					w,
					http.StatusInternalServerError,
					"Failed to get categories",
				)
				return
			}

			post.Categories = append(
				post.Categories,
				name,
			)
		}

		catRows.Close()
	}

	if post.Categories == nil {
		post.Categories = []string{}
	}

	// 5. Get comments

	commentRows, err := db.DB.Query(`
		SELECT
			c.id,
			u.nickname,
			c.content,
			c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)

	if err == nil {

		for commentRows.Next() {

			var c Comment

			if err := commentRows.Scan(
				&c.ID,
				&c.Author,
				&c.Content,
				&c.CreatedAt,
			); err != nil {

				commentRows.Close()

				helpers.SendJSON(
					w,
					http.StatusInternalServerError,
					"Failed to get comments",
				)
				return
			}

			post.Comments = append(
				post.Comments,
				c,
			)
		}

		commentRows.Close()
	}

	if post.Comments == nil {
		post.Comments = []Comment{}
	}

	// 6. Get total likes/dislikes

	stats, err := db.CountReactions(postID)

	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusInternalServerError,
			"Failed to get reaction stats",
		)
		return
	}

	// 7. Get current user's reaction

	userReaction, err := db.GetReaction(
		postID,
		userID,
	)

	if err != nil {
		helpers.SendJSON(
			w,
			http.StatusInternalServerError,
			"Failed to get user reaction",
		)
		return
	}

	// 8. Put reaction information into PostDetail

	post.Reactions.Likes = stats.Likes
	post.Reactions.Dislikes = stats.Dislikes

	if userReaction != nil {
		post.Reactions.UserReaction = userReaction.Reaction
	} else {
		post.Reactions.UserReaction = ""
	}

	// 9. Send JSON

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(post)
	fmt.Println("post", post)
}