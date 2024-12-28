package main

import (
	"fmt"

	"github.com/K1ender/MemeWhisper/internal/cache"
	"github.com/K1ender/MemeWhisper/internal/config"
	"github.com/K1ender/MemeWhisper/internal/database"
	"github.com/K1ender/MemeWhisper/internal/router"

	"go.uber.org/zap"
)

func main() {
	cfg := config.MustInit()

	var logger *zap.Logger

	if cfg.ENV == config.ProdENV {
		logger = zap.Must(zap.NewProduction())
	} else if cfg.ENV == config.LocalENV {
		logger = zap.Must(zap.NewDevelopment())
	}

	logger.Debug("Logger initialized")
	defer logger.Sync()

	logger.Debug("Connecting to database...")
	db := database.MustInit(cfg)
	defer db.Close()
	logger.Debug("Connected to database")

	logger.Debug("Connecting to memcached...")
	mc := cache.MustInit(cfg)
	defer mc.Close()
	logger.Debug("Connected to memcached")

	router := router.NewRouter(logger, db)

	logger.Debug("Starting server...")
	app := router.MustInit()
	logger.Debug("Server started")

	// userStore := store.NewUserStore(db, logger)
	// sessionStore := store.NewSessionStore(db, logger)

	logger.Error("Something went wrong",
		zap.Error(app.Listen(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))),
	)
}
