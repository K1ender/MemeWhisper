package store

import (
	"database/sql"

	"go.uber.org/zap"
)

type IMemeStore interface {
}

type memeStore struct {
	conn   *sql.DB
	logger *zap.Logger
}

func NewMemeStore(conn *sql.DB, logger *zap.Logger) IMemeStore {
	return &memeStore{
		conn: conn,
		logger: logger.With(
			zap.String("store", "user"),
		),
	}
}


