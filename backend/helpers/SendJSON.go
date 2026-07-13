package helpers

import (
	"encoding/json"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, message string, data ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(data) > 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": message, "data": data[0]})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": message})
	}
}
