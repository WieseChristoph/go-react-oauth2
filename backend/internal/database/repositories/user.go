package repositories

import (
	"github.com/WieseChristoph/go-oauth2-backend/internal/database"
	"github.com/WieseChristoph/go-oauth2-backend/internal/database/models"
)

type UserRepository interface {
	CreateUser(user *models.User) (int, error)
	UpdateUser(user *models.User) error
	GetAllUsers() ([]*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user *models.User) (int, error) {
	dbUser := user.ToDBUser()

	var id int
	err := r.db.QueryRow(`
		INSERT INTO user_data (name, display_name, email, avatar, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, dbUser.Name, dbUser.DisplayName, dbUser.Email, dbUser.Avatar, dbUser.Role).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	dbUser := user.ToDBUser()

	_, err := r.db.Exec(`
		UPDATE user_data
		SET name = $1, display_name = $2, email = $3, avatar = $4, role = $5
		WHERE id = $6
	`, dbUser.Name, dbUser.DisplayName, dbUser.Email, dbUser.Avatar, dbUser.Role, dbUser.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetAllUsers() ([]*models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, name, display_name, email, avatar, role, created_at, updated_at
		FROM user_data
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.Name, &user.DisplayName, &user.Email, &user.Avatar, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	user := models.User{}
	err := r.db.QueryRow(`
		SELECT id, name, display_name, email, avatar, role, created_at, updated_at
		FROM user_data
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.DisplayName, &user.Email, &user.Avatar, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}
	err := r.db.QueryRow(`
		SELECT id, name, display_name, email, avatar, role, created_at, updated_at
		FROM user_data
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.DisplayName, &user.Email, &user.Avatar, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
