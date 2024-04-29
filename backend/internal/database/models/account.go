package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrInvalidUserID            = errors.New("invalid user ID")
	ErrInvalidProvider          = errors.New("invalid provider")
	ErrInvalidProviderAccountID = errors.New("invalid provider account ID")
)

type APIAccount struct {
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

type Account struct {
	ID                int            `json:"id"`
	UserID            int            `json:"user_id"`
	Provider          string         `json:"provider"`
	ProviderAccountID string         `json:"provider_account_id"`
	AccessToken       sql.NullString `json:"access_token"`
	RefreshToken      sql.NullString `json:"refresh_token"`
	ExpiresAt         sql.NullTime   `json:"expires_at"`
	Scope             sql.NullString `json:"scope"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

func NewAccount(userID int, provider, providerAccountID, accessToken, refreshToken string, expiresAt time.Time, scope string) *Account {
	return &Account{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
		AccessToken:       sql.NullString{String: accessToken, Valid: accessToken != ""},
		RefreshToken:      sql.NullString{String: refreshToken, Valid: refreshToken != ""},
		ExpiresAt:         sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
		Scope:             sql.NullString{String: scope, Valid: scope != ""},
	}
}

func (a *Account) Validate() error {
	if a.UserID <= 0 {
		return ErrInvalidUserID
	}

	if a.Provider == "" {
		return ErrInvalidProvider
	}

	if a.ProviderAccountID == "" {
		return ErrInvalidProviderAccountID
	}

	return nil
}

func (a *Account) ToAPIAccount() *APIAccount {
	return &APIAccount{
		ID:                a.ID,
		UserID:            a.UserID,
		Provider:          a.Provider,
		ProviderAccountID: a.ProviderAccountID,
		AccessToken:       a.AccessToken.String,
		RefreshToken:      a.RefreshToken.String,
		ExpiresAt:         a.ExpiresAt.Time,
		Scope:             a.Scope.String,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

func (a *Account) UpdateTokens(accessToken, refreshToken string, expiresAt time.Time) {
	a.AccessToken = sql.NullString{String: accessToken, Valid: accessToken != ""}
	a.RefreshToken = sql.NullString{String: refreshToken, Valid: refreshToken != ""}
	a.ExpiresAt = sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()}
}

func (a *APIAccount) ToAccount() *Account {
	account := NewAccount(
		a.UserID,
		a.Provider,
		a.ProviderAccountID,
		a.AccessToken,
		a.RefreshToken,
		a.ExpiresAt,
		a.Scope,
	)
	account.ID = a.ID
	account.CreatedAt = a.CreatedAt
	account.UpdatedAt = a.UpdatedAt

	return account
}
