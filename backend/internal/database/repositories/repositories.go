package repositories

import (
	"github.com/WieseChristoph/go-oauth2-backend/internal/database"
)

type Repositories struct {
	User    UserRepository
	Account AccountRepository
	Session SessionRepository
}

func New(db *database.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Account: NewAccountRepository(db),
		Session: NewSessionRepository(db),
	}
}
