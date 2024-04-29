package api

import (
	"encoding/json"
	"net/http"

	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
)

func (a *api) GetUsersMe(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(config.UserContextKey).(*models.User)
	if !ok {
		log.Errorln("Error getting user from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user.ToAPIUser())
	if err != nil {
		log.Errorf("Error handling JSON marshal. Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(userJson)
}

func (a *api) GetUsersAll(w http.ResponseWriter, r *http.Request) {
	users, err := a.repositories.User.GetAllUsers()
	if err != nil {
		log.Errorf("Error getting users from database. Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	apiUsers := make([]*models.APIUser, 0, len(users))
	for _, user := range users {
		apiUsers = append(apiUsers, user.ToAPIUser())
	}

	usersJson, err := json.Marshal(apiUsers)
	if err != nil {
		log.Errorf("Error handling JSON marshal. Err: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(usersJson)
}
