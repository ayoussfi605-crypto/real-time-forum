package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	db "forum/database"
	"forum/helpers"
)

// Shape of a single comment
type Comment struct {
	ID        int    `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// Shape of a post with all its details (categories + comments included)
type PostDetail struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Author     string    `json:"author"`
	CreatedAt  string    `json:"created_at"`
	Categories []string  `json:"categories"`
	Comments   []Comment `json:"comments"`
}

// GET /api/posts/{id} — returns one post with its categories and comments
func GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {

	// r.PathValue("id") reads the {id} part from the URL
	// For example: /api/posts/5 → idStr = "5"
	idStr := r.PathValue("id")
	postID, err := strconv.Atoi(idStr) // convert "5" (string) to 5 (number)
	if err != nil {
		helpers.SendJSON(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post PostDetail

	// Get the post and the author's nickname
	err = db.DB.QueryRow(`
		SELECT p.id, p.title, p.content, p.created_at, u.nickname
		FROM posts p
		JOIN users u ON u.id = p.user_id
		WHERE p.id = ?
	`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Author)

	if err != nil {
		helpers.SendJSON(w, http.StatusNotFound, "Post not found")
		return
	}

	// Get the categories for this post
	catRows, err := db.DB.Query(`
		SELECT c.name FROM categories c
		JOIN post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = ?
	`, postID)
	if err == nil {
		for catRows.Next() {
			var name string
			catRows.Scan(&name)
			post.Categories = append(post.Categories, name)
		}
		catRows.Close() // explicitly close here before the next query
	}
	if post.Categories == nil {
		post.Categories = []string{}
	}

	// Get all comments for this post (oldest first)
	commentRows, err := db.DB.Query(`
		SELECT c.id, u.nickname, c.content, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err == nil {
		defer commentRows.Close()
		for commentRows.Next() {
			var c Comment
			commentRows.Scan(&c.ID, &c.Author, &c.Content, &c.CreatedAt)
			post.Comments = append(post.Comments, c)
		}
	}
	if post.Comments == nil {
		post.Comments = []Comment{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
