package helper

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func GetRandomElement(slice []string) string {
	rand.Seed(time.Now().UnixNano())

	if len(slice) == 0 {
		return ""
	}

	randomIndex := rand.Intn(len(slice))
	return slice[randomIndex]
}
