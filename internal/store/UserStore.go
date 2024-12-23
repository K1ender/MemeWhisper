package store

import (
	"database/sql"
	"errors"

	"github.com/K1ender/MemeWhisper/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserStore interface {
	GetUserByID(id int) (models.User, error)
	CreateUser(user models.User) error
	ChangeUsername(userID int, newUsername string) error
	ChangePassword(userID int, newPassword string) error
}

type UserStore struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewUserStore(conn *sql.DB, logger *zap.Logger) IUserStore {
	return &UserStore{
		conn: conn,
		logger: logger.With(
			zap.String("store", "user"),
		),
	}
}

func (s *UserStore) GetUserByID(id int) (models.User, error) {
	s.logger.Debug("Getting user by id", zap.Int("id", id))

	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return models.User{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, username, hashed_password FROM users WHERE id = $1", id)

	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.HashedPassword)

	if err != nil {
		s.logger.Error("Failed to scan user", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrUserDoesntExist
		}
		return models.User{}, ErrFailedToScanUser
	}

	err = tx.Commit()

	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.User{}, ErrFailedToCommitTransaction
	}

	s.logger.Debug("User found", zap.Any("userID", user.ID))
	return user, nil
}

func (s *UserStore) CreateUser(user models.User) error {
	s.logger.Debug("Creating user", zap.Any("userID", user.ID))
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)

	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return ErrFailedToHashPassword
	}

	_, err = tx.Exec("INSERT INTO users (username, hashed_password) VALUES ($1, $2)", user.Username, hashedPassword)

	if err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return ErrFailedToCreateUser
	}

	err = tx.Commit()

	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return ErrFailedToCommitTransaction
	}

	return nil
}

func (s *UserStore) ChangeUsername(userID int, newUsername string) error {
	s.logger.Debug(
		"Changing username",
		zap.Int("userID", userID),
		zap.String("newUsername", newUsername),
	)
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE users SET username = $1 WHERE id = $2", newUsername, userID)

	if err != nil {
		s.logger.Error("Failed to update username", zap.Error(err))
		return ErrFailedToUpdateUsername
	}

	err = tx.Commit()

	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return ErrFailedToCommitTransaction
	}

	return nil
}

func (s *UserStore) ChangePassword(userID int, newPassword string) error {
	s.logger.Debug(
		"Changing password",
		zap.Int("userID", userID),
	)
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return ErrFailedToHashPassword
	}

	_, err = tx.Exec("UPDATE users SET hashed_password = $1 WHERE id = $2", hashedPassword, userID)

	if err != nil {
		s.logger.Error("Failed to update password", zap.Error(err))
		return ErrFailedToUpdatePassword
	}

	err = tx.Commit()

	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return ErrFailedToCommitTransaction
	}

	return nil
}
