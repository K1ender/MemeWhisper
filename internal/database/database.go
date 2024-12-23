package database

import (
	"database/sql"
	"fmt"

	"github.com/K1ender/MemeWhisper/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func MustInit(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db
}
