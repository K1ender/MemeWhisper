package store

import (
	"database/sql"
	"errors"

	"github.com/K1ender/MemeWhisper/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserStore interface {
	GetUserByID(id int) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	CreateUser(user models.User) (int, error)
	ChangeUsername(userID int, newUsername string) error
	ChangePassword(userID int, newPassword string) error
}

type userStore struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewUserStore(conn *sql.DB, logger *zap.Logger) IUserStore {
	return &userStore{
		conn: conn,
		logger: logger.With(
			zap.String("store", "user"),
		),
	}
}

// GetUserByID retrieves a user from the database by their ID.
//
// This function takes the user's ID as an argument and returns a User model
// and an error. The error is non-nil if there was an error retrieving the user
// from the database.
func (s *userStore) GetUserByID(id int) (models.User, error) {
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

// GetUserByUsername retrieves a user from the database by their username.
//
// This function takes the user's username as an argument and returns a User model
// and an error. The error is non-nil if there was an error retrieving the user
// from the database.
func (s *userStore) GetUserByUsername(username string) (models.User, error) {
	s.logger.Debug("Getting user by username", zap.String("username", username))

	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return models.User{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id, username, hashed_password FROM users WHERE username = $1", username)

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

// CreateUser creates a new user in the database.
//
// This function takes a User model as an argument and returns the ID of the new
// user and an error. The error is non-nil if there was an error creating the user.
func (s *userStore) CreateUser(user models.User) (int, error) {
	s.logger.Debug("Creating user", zap.Any("userID", user.ID))
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return 0, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)

	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return 0, ErrFailedToHashPassword
	}

	row := tx.QueryRow("INSERT INTO users (username, hashed_password) VALUES ($1, $2) RETURNING id", user.Username, hashedPassword)

	if row.Err() != nil {
		var pgErr *pgconn.PgError
		if errors.As(row.Err(), &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, ErrUserAlreadyExists
		}

		s.logger.Error("Failed to create user", zap.Error(row.Err()))
		return 0, ErrFailedToCreateUser
	}

	var lastInsertID int
	err = row.Scan(&lastInsertID)

	if err != nil {
		s.logger.Error("Failed to get user id", zap.Error(err))
		return 0, ErrFailedToGetUserID
	}

	err = tx.Commit()

	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return 0, ErrFailedToCommitTransaction
	}

	return lastInsertID, nil
}

// ChangeUsername changes the username of the user with the given ID.
//
// This function takes a user ID and a new username as arguments and returns an
// error. The error is non-nil if there was an error changing the username.
func (s *userStore) ChangeUsername(userID int, newUsername string) error {
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

// ChangePassword changes the password of the user with the given ID.
//
// This function takes a user ID and a new password as arguments and returns an
// error. The error is non-nil if there was an error changing the password.
func (s *userStore) ChangePassword(userID int, newPassword string) error {
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
