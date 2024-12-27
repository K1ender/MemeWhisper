package router

import (
	"database/sql"

	"github.com/K1ender/MemeWhisper/internal/store"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Router struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewRouter(logger *zap.Logger, db *sql.DB) *Router {
	return &Router{
		logger: logger,
		db:     db,
	}
}

func (r *Router) MustInit() *fiber.App {
	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello, world!")
	})

	userStore := store.NewUserStore(r.db, r.logger.Named("user"))
	sessionStore := store.NewSessionStore(r.db, r.logger.Named("session"))

	userv1 := app.Group("/api/v1/user")
	NewUserRouter(
		userv1,
		r.logger.Named("user"),
		userStore,
	)

	sessionv1 := app.Group("/api/v1/session")
	NewSessionRouter(
		sessionv1,
		r.logger.Named("session"),
		userStore,
		sessionStore,
	)

	return app
}
