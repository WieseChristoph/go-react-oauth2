package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/auth"
	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/cookies"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
)

// Get the user with the session token and add it to the context
func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the cookie
		sessionToken, err := cookies.Read(r, config.SessionCookieName)
		if err != nil {
			if err != http.ErrNoCookie {
				log.Errorf("Error reading session cookie. Err: %v", err)
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get the session from the database
		session, err := m.repositories.Session.GetSessionByToken(sessionToken)
		if err != nil {
			if err == sql.ErrNoRows {
				// Delete the session cookie if the session does not exist
				cookies.Delete(w, config.SessionCookieName)
			} else {
				log.Errorf("Error getting session from database. Err: %v", err)
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the session is expired
		if session.IsExpired() {
			// Delete the session cookie and the session from the database
			cookies.Delete(w, config.SessionCookieName)
			err = m.repositories.Session.DeleteSessionByToken(sessionToken)
			if err != nil {
				log.Errorf("Error deleting session from database. Err: %v", err)
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the session needs to be refreshed
		if session.IsRefreshNeeded() {
			// Create a new session cookie with the same token to refresh the expiry
			sessionCookie := auth.NewSessionCookie(sessionToken)
			err := cookies.Write(w, sessionCookie)
			if err != nil {
				log.Errorf("Error writing session cookie. Err: %v", err)
			} else {
				// Update the expiry of the session in the database
				session.ExpiresAt = time.Now().Add(config.SessionMaxAge)
				err = m.repositories.Session.UpdateSessionExpiresAtByToken(sessionToken, session.ExpiresAt)
				if err != nil {
					log.Errorf("Error updating expiry of session in database. Err: %v", err)
				}
			}
		}

		// Get the user from the session
		user, err := m.repositories.User.GetUserByID(session.UserID)
		if err != nil {
			log.Errorf("Error getting user from database. Err: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the user to the context
		ctx := context.WithValue(r.Context(), config.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to check if the user has any of the required roles
func (m *middleware) RoleMiddleware(roles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user from the context
			user, ok := r.Context().Value(config.UserContextKey).(*models.User)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the user has any of the required roles
			allowed := false
			for _, role := range roles {
				if user.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
