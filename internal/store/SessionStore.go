package store

import (
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"strings"

	"go.uber.org/zap"
)

type ISessionStore interface {
	generateSessionToken() (string, error)
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

func (s *sessionStore) generateSessionToken() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Error("Failed to generate random bytes", zap.Error(err))
		return "", ErrFailedToGenerateRandomBytes
	}
	return strings.ToLower(base32.HexEncoding.EncodeToString(bytes)), nil
}
