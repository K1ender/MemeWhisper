package store

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/hex"
	"strings"
	"time"

	"github.com/K1ender/MemeWhisper/internal/models"
	"go.uber.org/zap"
)

type ISessionStore interface {
	GenerateSessionToken() (string, error)
	CreateSession(token string, userID int) (models.Session, error)
	ValidateSessionToken(token string) (models.Session, error)
	InvalidateSession(sessionID string) error
}

type sessionStore struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewSessionStore(conn *sql.DB, logger *zap.Logger) ISessionStore {
	return &sessionStore{
		conn: conn,
		logger: logger.With(
			zap.String("store", "session"),
		),
	}
}

// GenerateSessionToken generates a random session token.
//
// This function returns a string and an error. The error is non-nil if there
// was an error generating the session token.
func (s *sessionStore) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Error("Failed to generate random bytes", zap.Error(err))
		return "", ErrFailedToGenerateRandomBytes
	}
	s.logger.Debug("Generated random bytes", zap.ByteString("bytes", bytes))
	token := strings.ToLower(base32.StdEncoding.EncodeToString(bytes))
	s.logger.Debug("Generated session token", zap.String("token", token))
	return token, nil
}

// CreateSession creates a new session in the database.
//
// This function takes a session token and a user ID as arguments and returns a Session model and an error.
// The error is non-nil if there was an error creating the session.
func (s *sessionStore) CreateSession(token string, userID int) (models.Session, error) {
	hash := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hash[:])
	session := models.Session{
		ID:        sessionId,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	s.logger.Debug("Creating session", zap.Any("session", session))

	tx, err := s.conn.Begin()

	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return models.Session{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	s.logger.Debug("Executing query", zap.String("query", "INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)"), zap.Any("args", []interface{}{sessionId, userID, session.ExpiresAt}))
	_, err = tx.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)", sessionId, userID, session.ExpiresAt)
	if err != nil {
		s.logger.Error("Failed to create session", zap.Error(err))
		return models.Session{}, ErrFailedToCreateSession
	}

	s.logger.Debug("Committing transaction")

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.Session{}, ErrFailedToCommitTransaction
	}

	s.logger.Debug("Session created successfully")

	return session, nil
}

// ValidateSessionToken validates a session token.
//
// This function takes a session token as an argument and returns a Session model
// and an error. The error is non-nil if there was an error validating the session
// token.
func (s *sessionStore) ValidateSessionToken(token string) (models.Session, error) {
	hash := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hash[:])

	s.logger.Debug("Validating session", zap.String("token", token), zap.String("sessionId", sessionId))

	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return models.Session{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	var session models.Session
	err = tx.QueryRow("SELECT id, user_id, expires_at FROM sessions WHERE id = $1", sessionId).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		s.logger.Error("Failed to scan session", zap.Error(err))
		return models.Session{}, ErrFailedToScanSession
	}

	s.logger.Debug("Session found", zap.Any("session", session))

	if time.Now().After(session.ExpiresAt) {
		s.logger.Debug("Session expired", zap.Any("session", session))
		tx.Exec("DELETE FROM sessions WHERE id = $1", sessionId)
		return models.Session{}, ErrSessionExpired
	}

	if time.Now().After(session.ExpiresAt.Add(-15 * 24 * time.Hour)) {
		session.ExpiresAt = time.Now().Add(time.Hour * 24 * 30)
		s.logger.Debug("Updating session", zap.Any("session", session))
		tx.Exec("UPDATE sessions SET expires_at = $1 WHERE id = $2", session.ExpiresAt, sessionId)
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.Session{}, ErrFailedToCommitTransaction
	}

	s.logger.Debug("Session validated", zap.Any("session", session))

	return session, nil
}

// InvalidateSession deletes a session from the database.
//
// This function takes a session ID as an argument and returns an error. The
// error is non-nil if there was an error deleting the session from the database.
func (s *sessionStore) InvalidateSession(sessionID string) error {
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return ErrFailedToStartTransaction
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("Failed to rollback transaction", zap.Error(err))
		}
	}()

	s.logger.Debug("Attempting to delete session", zap.String("sessionID", sessionID))
	_, err = tx.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		s.logger.Error("Failed to delete session", zap.Error(err), zap.String("sessionID", sessionID))
		return ErrFailedToDeleteSession
	}
	s.logger.Debug("Session deleted successfully", zap.String("sessionID", sessionID))

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return ErrFailedToCommitTransaction
	}
	s.logger.Debug("Transaction committed successfully")

	return nil
}
