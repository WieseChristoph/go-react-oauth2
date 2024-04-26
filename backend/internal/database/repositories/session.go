package repositories

import (
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/database"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
)

type SessionRepository interface {
	CreateSession(session *models.Session) (int, error)
	UpdateSessionExpiresAtByToken(token string, expiresAt time.Time) error
	DeleteSessionByToken(token string) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteExpiredSessions() error
}

type sessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) SessionRepository {
	return &sessionRepository{
		db: db,
	}
}

func (r *sessionRepository) CreateSession(session *models.Session) (int, error) {
	dbSession := session.ToDBSession()

	var id int
	err := r.db.QueryRow(`
		INSERT INTO session (token, user_id, ip_address, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, dbSession.Token, dbSession.UserID, dbSession.IPAddress, dbSession.ExpiresAt).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *sessionRepository) UpdateSessionExpiresAtByToken(token string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		UPDATE session
		SET expires_at = $1
		WHERE token = $2
	`, expiresAt, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) DeleteSessionByToken(token string) error {
	_, err := r.db.Exec(`
		DELETE FROM session
		WHERE token = $1
	`, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) GetSessionByToken(token string) (*models.Session, error) {
	session := models.Session{}
	err := r.db.QueryRow(`
		SELECT id, user_id, token, ip_address, expires_at, created_at, updated_at
		FROM session
		WHERE token = $1
	`, token).Scan(&session.ID, &session.UserID, &session.Token, &session.IPAddress, &session.ExpiresAt, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteExpiredSessions() error {
	_, err := r.db.Exec(`
		DELETE FROM session
		WHERE expires_at < NOW()
	`)
	if err != nil {
		return err
	}

	return nil
}
