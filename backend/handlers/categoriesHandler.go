package handlers

import (
	"encoding/json"
	"net/http"

	db "forum/database"
	"forum/helpers"
)

// Shape of a category
type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GET /api/categories — returns all categories
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Could not get categories")
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			helpers.SendJSON(w, http.StatusInternalServerError, "Error reading categories")
			return
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		helpers.SendJSON(w, http.StatusInternalServerError, "Error reading categories")
		return
	}

	if categories == nil {
		categories = []Category{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
