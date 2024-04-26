package api

import (
	"net/http"

	"github.com/WieseChristoph/go-oauth2-backend/internal/database/repositories"
)

type Api interface {
	GetUsersMe(w http.ResponseWriter, r *http.Request)
	GetUsersAll(w http.ResponseWriter, r *http.Request)
}

type api struct {
	repositories *repositories.Repositories
}

func New(repositories *repositories.Repositories) Api {
	return &api{
		repositories: repositories,
	}
}
