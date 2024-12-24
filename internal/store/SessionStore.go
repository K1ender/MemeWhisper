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

func (s *sessionStore) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Error("Failed to generate random bytes", zap.Error(err))
		return "", ErrFailedToGenerateRandomBytes
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(bytes)), nil
}

func (s *sessionStore) CreateSession(token string, userID int) (models.Session, error) {
	hash := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hash[:])
	session := models.Session{
		ID:        sessionId,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	tx, err := s.conn.Begin()

	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return models.Session{}, ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)", sessionId, userID, session.ExpiresAt)
	if err != nil {
		s.logger.Error("Failed to create session", zap.Error(err))
		return models.Session{}, ErrFailedToCreateSession
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.Session{}, ErrFailedToCommitTransaction
	}

	return session, nil
}

func (s *sessionStore) ValidateSessionToken(token string) (models.Session, error) {
	hash := sha256.Sum256([]byte(token))
	sessionId := hex.EncodeToString(hash[:])

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

	if time.Now().After(session.ExpiresAt) {
		tx.Exec("DELETE FROM sessions WHERE id = $1", sessionId)
		return models.Session{}, ErrSessionExpired
	}

	if time.Now().After(session.ExpiresAt.Add(-15 * 24 * time.Hour)) {
		session.ExpiresAt = time.Now().Add(time.Hour * 24 * 30)
		tx.Exec("UPDATE sessions SET expires_at = $1 WHERE id = $2", session.ExpiresAt, sessionId)
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return models.Session{}, ErrFailedToCommitTransaction
	}

	return session, nil
}

func (s *sessionStore) InvalidateSession(sessionID string) error {
	tx, err := s.conn.Begin()
	if err != nil {
		s.logger.Error("Failed to start transaction", zap.Error(err))
		return ErrFailedToStartTransaction
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		s.logger.Error("Failed to delete session", zap.Error(err))
		return ErrFailedToDeleteSession
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return ErrFailedToCommitTransaction
	}

	return nil
}
