package models

import (
	"database/sql"
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
)

type Session struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	UserID    int       `json:"user_id"`
	IPAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBSession struct {
	ID        int            `json:"id"`
	Token     sql.NullString `json:"token"`
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
		IPAddress: ipAddress,
		ExpiresAt: expiresAt,
	}
}

func (s *Session) ToDBSession() *DBSession {
	return &DBSession{
		ID:        s.ID,
		Token:     sql.NullString{String: s.Token, Valid: s.Token != ""},
		UserID:    s.UserID,
		IPAddress: sql.NullString{String: s.IPAddress, Valid: s.IPAddress != ""},
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) IsRefreshNeeded() bool {
	return s.ExpiresAt.Add(-config.SessionMaxAge + config.SessionRefreshAge).Before(time.Now())
}
