package models

import (
	"database/sql"
	"time"
)

type Account struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	Provider          string    `json:"provider"`
	ProviderAccountID string    `json:"provider_account_id"`
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token"`
	ExpiresAt         time.Time `json:"expires_at"`
	Scope             string    `json:"scope"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type DBAccount struct {
	ID                int            `json:"id"`
	UserID            int            `json:"user_id"`
	Provider          sql.NullString `json:"provider"`
	ProviderAccountID sql.NullString `json:"provider_account_id"`
	AccessToken       sql.NullString `json:"access_token"`
	RefreshToken      sql.NullString `json:"refresh_token"`
	ExpiresAt         time.Time      `json:"expires_at"`
	Scope             sql.NullString `json:"scope"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func NewAccount(userID int, provider, providerAccountID, accessToken, refreshToken string, expiresAt time.Time, scope string) *Account {
	return &Account{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
		ExpiresAt:         expiresAt,
		Scope:             scope,
	}
}

func (a *Account) ToDBAccount() *DBAccount {
	return &DBAccount{
		ID:                a.ID,
		UserID:            a.UserID,
		Provider:          sql.NullString{String: a.Provider, Valid: a.Provider != ""},
		ProviderAccountID: sql.NullString{String: a.ProviderAccountID, Valid: a.ProviderAccountID != ""},
		AccessToken:       sql.NullString{String: a.AccessToken, Valid: a.AccessToken != ""},
		RefreshToken:      sql.NullString{String: a.RefreshToken, Valid: a.RefreshToken != ""},
		ExpiresAt:         a.ExpiresAt,
		Scope:             sql.NullString{String: a.Scope, Valid: a.Scope != ""},
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}
