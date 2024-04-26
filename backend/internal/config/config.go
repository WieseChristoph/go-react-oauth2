package config

import "time"

type ContextKey string

const (
	SessionCookieName = "session"
	SessionMaxAge     = time.Hour * 24 * 30
	SessionRefreshAge = time.Hour * 24 * 7

	UserContextKey     ContextKey = "user"
	ProviderContextKey ContextKey = "provider"
)
