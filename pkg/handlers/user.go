package handlers

import (
	"encoding/json"
	"net/http"

	_ "golang.org/x/crypto/bcrypt"
	_ "video-conference/pkg/models"
	_ "video-conference/pkg/utils"
)

func (repo *Repos) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	users, err := repo.UserRepo.FindUsersByUsername(query)
	if err != nil {
		http.Error(w, "Error in fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
