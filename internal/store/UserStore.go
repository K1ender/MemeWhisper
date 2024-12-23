package store

import (
	"database/sql"
	"errors"

	"github.com/K1ender/MemeWhisper/internal/models"
)

type IUserStore interface {
	GetUserByID(id int) (models.User, error)
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
}

type UserStore struct {
	conn *sql.DB
}

func NewUserStore(conn *sql.DB) IUserStore {
	return &UserStore{
		conn: conn,
	}
}

func (s *UserStore) GetUserByID(id int) (models.User, error) {

	tx, err := s.conn.Begin()
	if err != nil {
		return models.User{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, username, hashed_password FROM users WHERE id = $1", id)

	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.HashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserDoesntExist
		}
		return models.User{}, ErrFailedToScanUser
	}

	err = tx.Commit()

	if err != nil {
		return models.User{}, ErrFailedToCommitTransaction
	}

	return user, nil
}

func (s *UserStore) CreateUser(user models.User) error {
	return nil
}

func (s *UserStore) UpdateUser(user models.User) error {
	return nil
}
