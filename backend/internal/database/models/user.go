package models

import (
	"database/sql"
	"time"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DBUser struct {
	ID          int            `json:"id"`
	Name        sql.NullString `json:"name"`
	DisplayName sql.NullString `json:"display_name"`
	Email       sql.NullString `json:"email"`
	Avatar      sql.NullString `json:"avatar"`
	Role        sql.NullString `json:"role"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func NewUser(name, displayName, email, avatar string, role Role) *User {
	return &User{
		Name:        name,
		DisplayName: displayName,
		Email:       email,
		Avatar:      avatar,
		Role:        role,
	}
}

func (u *User) ToDBUser() *DBUser {
	if u.Role == "" {
		u.Role = RoleUser
	}

	return &DBUser{
		ID:          u.ID,
		Name:        sql.NullString{String: u.Name, Valid: u.Name != ""},
		DisplayName: sql.NullString{String: u.DisplayName, Valid: u.DisplayName != ""},
		Email:       sql.NullString{String: u.Email, Valid: u.Email != ""},
		Avatar:      sql.NullString{String: u.Avatar, Valid: u.Avatar != ""},
		Role:        sql.NullString{String: string(u.Role), Valid: u.Role != ""},
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
