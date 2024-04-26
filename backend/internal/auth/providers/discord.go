package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
	"golang.org/x/oauth2"
)

const (
	discordProviderName = "discord"
	discordAuthURL      = "https://discord.com/api/oauth2/authorize"
	discordTokenURL     = "https://discord.com/api/oauth2/token"
	discordUserInfoURL  = "https://discord.com/api/users/@me"
	discordAvatarURL    = "https://cdn.discordapp.com/avatars/%s/%s.png"
)

type discordUserInfo struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email"`
	Verified      bool   `json:"verified"`
}

type DiscordProvider struct {
	name   string
	config *oauth2.Config
}

func NewDiscordProvider(clientID string, clientSecret string, redirectURL string, scopes []string) OAuth2Provider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  discordAuthURL,
			TokenURL: discordTokenURL,
		},
	}

	return &DiscordProvider{
		name:   discordProviderName,
		config: config,
	}
}

func (p *DiscordProvider) Name() string {
	return p.name
}

func (p *DiscordProvider) Config() *oauth2.Config {
	return p.config
}

func (p *DiscordProvider) UserInfo(token *oauth2.Token) (*models.User, string, error) {
	res, err := p.config.Client(context.Background(), token).Get(discordUserInfoURL)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	var userInfo discordUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, "", err
	}

	user := models.NewUser(
		userInfo.Username,
		userInfo.GlobalName,
		userInfo.Email,
		fmt.Sprintf(discordAvatarURL, userInfo.ID, userInfo.Avatar),
		"",
	)

	return user, userInfo.ID, nil
}
