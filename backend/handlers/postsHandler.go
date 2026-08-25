package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	db "forum/database"
	"forum/helpers"
	"forum/middlewares"
)

// This is the shape of what we send back for each post
type Post struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Author     string   `json:"author"`
	CreatedAt  string   `json:"created_at"`
	Categories []string `json:"categories"`
}

// This is what we expect from the frontend when creating a post
type CreatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Categories []int  `json:"categories"` // list of category IDs (e.g. [1, 3])
}

// GET /api/posts — returns all posts (newest first)
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {

	// Ask the database for all posts, and also grab the author's nickname and categories using GROUP_CONCAT
	rows, err := db.DB.Query(`
		SELECT 
			p.id, 
			p.title, 
			p.content, 
			p.created_at, 
			u.nickname,
			IFNULL(GROUP_CONCAT(c.name), '') as categories
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN post_categories pc ON pc.post_id = p.id
		LEFT JOIN categories c ON c.id = pc.category_id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not get posts")
		return
	}
	defer rows.Close()

	var posts []Post

	// Loop through each row (each post) returned by the database
	for rows.Next() {
		var p Post
		var categoriesStr string

		// Fill in the post fields from the database row
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Author, &categoriesStr)
		if err != nil {
			helpers.SendJSON(w, http.StatusInternalServerError, "Error reading posts")
			return
		}

		if categoriesStr != "" {
			p.Categories = strings.Split(categoriesStr, ",")
		} else {
			p.Categories = []string{}
		}

		posts = append(posts, p)
	}

	// Make sure posts is an empty list [] instead of null in JSON
	if posts == nil {
		posts = []Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

// POST /api/posts — creates a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SendJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(int)
	if !ok {
		helpers.SendJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		helpers.SendJSON(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validation
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)

	if input.Title == "" || input.Content == "" {
		helpers.SendJSON(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	if len(input.Categories) == 0 {
		helpers.SendJSON(w, http.StatusBadRequest, "Please select at least one category")
		return
	}

	if len(input.Title) > 100 {
		helpers.SendJSON(w, http.StatusBadRequest, "Title is too long")
		return
	}

	if len(input.Content) > 10000 {
		helpers.SendJSON(w, http.StatusBadRequest, "Content is too long")
		return
	}

	// Check that all categories exist BEFORE creating the post
	for _, catID := range input.Categories {
		var exists int

		err := db.DB.QueryRow(
			"SELECT 1 FROM categories WHERE id = ?",
			catID,
		).Scan(&exists)

		if err != nil {
			helpers.SendJSON(
				w,
				http.StatusBadRequest,
				"Category does not exist",
			)
			return
		}
	}

	// Start transaction
	tx, err := db.DB.Begin()
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not start transaction")
		return
	}

	// Create post
	result, err := tx.Exec(
		"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
		userID,
		input.Title,
		input.Content,
	)
	if err != nil {
		tx.Rollback()
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not create post")
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not get post ID")
		return
	}

	// Link categories
	for _, catID := range input.Categories {
		_, err := tx.Exec(
			"INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
			postID,
			catID,
		)

		if err != nil {
			tx.Rollback()
			helpers.SendJSON(
				w,
				http.StatusInternalServerError,
				"Could not assign category",
			)
			return
		}
	}

	// Everything succeeded
	if err := tx.Commit(); err != nil {
		helpers.SendJSON(
			w,
			http.StatusInternalServerError,
			"Could not save post",
		)
		return
	}

	helpers.SendJSON(w, http.StatusCreated, "Post created successfully")
}