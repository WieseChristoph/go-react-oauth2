package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidExpiresAt = errors.New("invalid expires at")
)

type APISession struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID        int            `json:"id"`
	Token     string         `json:"token"`
	UserID    int            `json:"user_id"`
	IPAddress sql.NullString `json:"ip_address"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func NewSession(token string, userID int, ipAddress string, expiresAt time.Time) *Session {
	return &Session{
		Token:     token,
		UserID:    userID,
		IPAddress: sql.NullString{String: ipAddress, Valid: ipAddress != ""},
		ExpiresAt: expiresAt,
	}
}

func (s *Session) Validate() error {
	if s.Token == "" {
		return ErrInvalidToken
	}

	if s.UserID <= 0 {
		return ErrInvalidUserID
	}

	if s.ExpiresAt.IsZero() {
		return ErrInvalidExpiresAt
	}

	return nil
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) IsRefreshNeeded() bool {
	return s.ExpiresAt.Add(-config.SessionMaxAge + config.SessionRefreshAge).Before(time.Now())
}

func (s *Session) ToAPISession() *APISession {
	return &APISession{
		ID:        s.ID,
		Token:     s.Token,
		UserID:    s.UserID,
		IPAddress: s.IPAddress.String,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (s *APISession) ToSession() *Session {
	session := NewSession(s.Token, s.UserID, s.IPAddress, s.ExpiresAt)
	session.ID = s.ID
	session.CreatedAt = s.CreatedAt
	session.UpdatedAt = s.UpdatedAt

	return session
}
