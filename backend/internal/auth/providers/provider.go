package providers

import (
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"golang.org/x/oauth2"
)

type OAuth2Provider interface {
	Name() string
	Config() *oauth2.Config
	UserInfo(token *oauth2.Token) (*models.User, string, error)
}
