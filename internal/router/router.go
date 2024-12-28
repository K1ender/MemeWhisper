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
	r.logger.Debug("Initializing Fiber app")

	app := fiber.New(fiber.Config{
		// Prefork: true,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		r.logger.Debug("Handling health check request")
		return c.Status(200).SendString("Hello, world!")
	})

	r.logger.Debug("Creating user store")
	userStore := store.NewUserStore(r.db, r.logger.Named("user"))

	r.logger.Debug("Creating session store")
	sessionStore := store.NewSessionStore(r.db, r.logger.Named("session"))

	userv1 := app.Group("/api/v1/user")
	r.logger.Debug("Creating user routes")
	NewUserRouter(
		userv1,
		r.logger.Named("user"),
		userStore,
	)

	sessionv1 := app.Group("/api/v1/session")
	r.logger.Debug("Creating session routes")
	NewSessionRouter(
		sessionv1,
		r.logger.Named("session"),
		userStore,
		sessionStore,
	)

	r.logger.Debug("Finished initializing Fiber app")
	return app
}
