package middleware

import (
	"net/http"

	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/repositories"
)

type Middleware interface {
	AuthMiddleware(next http.Handler) http.Handler
	RoleMiddleware(roles ...models.Role) func(http.Handler) http.Handler
}

type middleware struct {
	repositories *repositories.Repositories
}

func New(repositories *repositories.Repositories) Middleware {
	return &middleware{repositories: repositories}
}
