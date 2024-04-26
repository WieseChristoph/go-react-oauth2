package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/auth/providers"
	"github.com/WieseChristoph/go-oauth2-backend/internal/config"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/repositories"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/cookies"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/env"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var (
	ErrProviderNotFound   = errors.New("provider not found")
	ErrOAuth2CodeNotFound = errors.New("oauth2 code not found")
	ErrAlreadyLoggedIn    = errors.New("already logged in")
)

type Auth interface {
	BeginOAuth2Flow(w http.ResponseWriter, r *http.Request) error
	HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) error
	Logout(w http.ResponseWriter, r *http.Request) error
}

type auth struct {
	providers    map[string]providers.OAuth2Provider
	repositories *repositories.Repositories
}

func New(repositories *repositories.Repositories) Auth {
	appURL := env.GetEnvOrFatal("APP_URL")

	authService := &auth{
		providers:    make(map[string]providers.OAuth2Provider),
		repositories: repositories,
	}

	// Add providers
	discordProvider := providers.NewDiscordProvider(
		env.GetEnvOrFatal("DISCORD_CLIENT_ID"),
		env.GetEnvOrFatal("DISCORD_CLIENT_SECRET"),
		fmt.Sprintf("%s/auth/discord/callback", appURL),
		[]string{"identify", "email"},
	)
	googleProvider := providers.NewGoogleProvider(
		env.GetEnvOrFatal("GOOGLE_CLIENT_ID"),
		env.GetEnvOrFatal("GOOGLE_CLIENT_SECRET"),
		fmt.Sprintf("%s/auth/google/callback", appURL),
		[]string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	)
	authService.addProviders(discordProvider, googleProvider)

	return authService
}

func (a *auth) BeginOAuth2Flow(w http.ResponseWriter, r *http.Request) error {
	// Check if user is already logged in
	sessionToken, err := cookies.Read(r, config.SessionCookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			return err
		}
	} else {
		_, err := a.repositories.Session.GetSessionByToken(sessionToken)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			// Delete invalid session cookie
			cookies.Delete(w, config.SessionCookieName)
		} else {
			return ErrAlreadyLoggedIn
		}
	}

	provider := r.Context().Value(config.ProviderContextKey).(string)

	state, err := generateToken()
	if err != nil {
		return err
	}

	p := a.providers[provider]
	if p == nil {
		return ErrProviderNotFound
	}

	u := p.Config().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	return nil
}

func (a *auth) HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) error {
	provider := r.Context().Value(config.ProviderContextKey).(string)

	code := r.FormValue("code")
	if code == "" {
		return ErrOAuth2CodeNotFound
	}

	p := a.providers[provider]
	if p == nil {
		return ErrProviderNotFound
	}

	// Exchange code for token
	token, err := p.Config().Exchange(r.Context(), code)
	if err != nil {
		return err
	}

	// Get user info
	user, providerAccountID, err := p.UserInfo(token)
	if err != nil {
		return err
	}

	// Check if account already exists
	// If account does not exist, create account and user
	// If account exists, update account and user
	account, err := a.repositories.Account.GetAccountByProviderAndProviderAccountID(provider, providerAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Create user
			user.ID, err = a.repositories.User.CreateUser(user)
			if err != nil {
				return err
			}

			// Create account
			account = models.NewAccount(
				user.ID,
				provider,
				providerAccountID,
				token.AccessToken,
				token.RefreshToken,
				token.Expiry,
				strings.Join(p.Config().Scopes, " "),
			)
			_, err = a.repositories.Account.CreateAccount(account)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Update account
		account.AccessToken = token.AccessToken
		account.RefreshToken = token.RefreshToken
		account.ExpiresAt = token.Expiry
		err = a.repositories.Account.UpdateAccount(account)
		if err != nil {
			return err
		}

		// Get user
		dbUser, err := a.repositories.User.GetUserByID(account.UserID)
		if err != nil {
			return err
		}

		// Update user
		user.ID = dbUser.ID
		user.Role = dbUser.Role
		err = a.repositories.User.UpdateUser(user)
		if err != nil {
			return err
		}
	}

	sessionToken, err := generateToken()
	if err != nil {
		return err
	}

	sessionCookie := NewSessionCookie(sessionToken)

	err = cookies.Write(w, sessionCookie)
	if err != nil {
		return err
	}

	// Format IP address
	ipAddress := r.RemoteAddr
	if lastIndex := strings.LastIndex(ipAddress, ":"); lastIndex != -1 {
		ipAddress = ipAddress[0:lastIndex]
	}
	ipAddress = strings.Trim(ipAddress, "[]")

	// Create session
	session := models.NewSession(
		sessionToken,
		user.ID,
		ipAddress,
		time.Now().Add(config.SessionMaxAge),
	)
	_, err = a.repositories.Session.CreateSession(session)
	if err != nil {
		return err
	}

	return nil
}

func (a *auth) Logout(w http.ResponseWriter, r *http.Request) error {
	sessionToken, err := cookies.Read(r, config.SessionCookieName)
	if err != nil {
		return err
	}

	err = a.repositories.Session.DeleteSessionByToken(sessionToken)
	if err != nil {
		return err
	}

	cookies.Delete(w, config.SessionCookieName)

	return nil
}

func (a *auth) addProviders(providers ...providers.OAuth2Provider) {
	for _, p := range providers {
		a.providers[p.Name()] = p
	}
}

func NewSessionCookie(sessionToken string) http.Cookie {
	appURL := env.GetEnvOrFatal("APP_URL")

	return http.Cookie{
		Name:     config.SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   int(config.SessionMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   strings.Contains(appURL, "https"),
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
