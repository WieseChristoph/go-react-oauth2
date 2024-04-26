package providers

import (
	"context"
	"encoding/json"
	"io"

	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"golang.org/x/oauth2"
)

const (
	googleProviderName  = "google"
	googleAuthURL       = "https://accounts.google.com/o/oauth2/auth"
	googleTokenURL      = "https://oauth2.googleapis.com/token"
	googleDeviceAuthURL = "https://oauth2.googleapis.com/device/code"
	googleUserInfoURL   = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
)

type googleUserInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Locale        string `json:"locale"`
}

type GoogleProvider struct {
	name   string
	config *oauth2.Config
}

func NewGoogleProvider(clientID string, clientSecret string, redirectURL string, scopes []string) OAuth2Provider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       googleAuthURL,
			TokenURL:      googleTokenURL,
			DeviceAuthURL: googleDeviceAuthURL,
		},
	}

	return &GoogleProvider{
		name:   googleProviderName,
		config: config,
	}
}

func (p *GoogleProvider) Name() string {
	return p.name
}

func (p *GoogleProvider) Config() *oauth2.Config {
	return p.config
}

func (p *GoogleProvider) UserInfo(token *oauth2.Token) (*models.User, string, error) {
	res, err := p.config.Client(context.Background(), token).Get(googleUserInfoURL)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	var userInfo googleUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, "", err
	}

	user := models.NewUser(
		userInfo.Name,
		userInfo.GivenName,
		userInfo.Email,
		userInfo.Picture,
		"",
	)

	return user, userInfo.ID, nil
}
