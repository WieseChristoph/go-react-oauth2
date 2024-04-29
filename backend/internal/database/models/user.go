package models

import (
	"database/sql"
	"errors"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidRole     = errors.New("invalid role")
)

type APIUser struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	DisplayName sql.NullString `json:"display_name"`
	Email       sql.NullString `json:"email"`
	Avatar      sql.NullString `json:"avatar"`
	Role        Role           `json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func NewUser(name, displayName, email, avatar string, role Role) *User {
	if role == "" {
		role = RoleUser
	}

	return &User{
		Name:        name,
		DisplayName: sql.NullString{String: displayName, Valid: displayName != ""},
		Email:       sql.NullString{String: email, Valid: email != ""},
		Avatar:      sql.NullString{String: avatar, Valid: avatar != ""},
		Role:        role,
	}
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}

	if u.Role != RoleUser && u.Role != RoleAdmin {
		return ErrInvalidRole
	}

	return nil
}

func (u *User) ToAPIUser() *APIUser {
	return &APIUser{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName.String,
		Email:       u.Email.String,
		Avatar:      u.Avatar.String,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func (u *APIUser) ToUser() *User {
	user := NewUser(
		u.Name,
		u.DisplayName,
		u.Email,
		u.Avatar,
		u.Role,
	)
	user.ID = u.ID
	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt

	return user
}
