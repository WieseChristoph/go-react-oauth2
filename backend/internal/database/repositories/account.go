package repositories

import (
	"github.com/WieseChristoph/go-oauth2-backend/internal/database"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
)

type AccountRepository interface {
	CreateAccount(account *models.Account) (int, error)
	UpdateAccount(account *models.Account) error
	GetAccountByProviderAndProviderAccountID(provider string, providerAccountID string) (*models.Account, error)
}

type accountRepository struct {
	db *database.DB
}

func NewAccountRepository(db *database.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) CreateAccount(account *models.Account) (int, error) {
	dbAccount := account.ToDBAccount()

	var id int
	err := r.db.QueryRow(`
		INSERT INTO account (user_id, provider, provider_account_id, access_token, refresh_token, expires_at, scope)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, dbAccount.UserID, dbAccount.Provider, dbAccount.ProviderAccountID, dbAccount.AccessToken, dbAccount.RefreshToken, dbAccount.ExpiresAt, dbAccount.Scope).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *accountRepository) UpdateAccount(account *models.Account) error {
	dbAccount := account.ToDBAccount()

	_, err := r.db.Exec(`
		UPDATE account
		SET access_token = $1, refresh_token = $2, expires_at = $3
		WHERE id = $4
	`, dbAccount.AccessToken, dbAccount.RefreshToken, dbAccount.ExpiresAt, dbAccount.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *accountRepository) GetAccountByProviderAndProviderAccountID(provider string, providerAccountID string) (*models.Account, error) {
	account := models.Account{}
	err := r.db.QueryRow(`
		SELECT id, user_id, provider, provider_account_id, access_token, refresh_token, expires_at, scope, created_at, updated_at
		FROM account
		WHERE provider = $1 AND provider_account_id = $2
	`, provider, providerAccountID).Scan(&account.ID, &account.UserID, &account.Provider, &account.ProviderAccountID, &account.AccessToken, &account.RefreshToken, &account.ExpiresAt, &account.Scope, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
